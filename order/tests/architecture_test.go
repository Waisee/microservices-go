package tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/waisee/microservices-go/order/tests/testutil"
)

// Тесты этого файла подчёркивают, что в week_2 появился слой бизнес-логики
// (service-слой Clean Architecture) и доменные ошибки. В week_1 не было такого
// разделения: HTTP-хендлер сам вычислял цену и работал со store.
//
// Здесь мы валидируем не отдельные эндпоинты, а свойство архитектуры:
//  1. Бизнес-вычисления (TotalPrice) делает service-слой, и они одинаковы
//     для любых комбинаций деталей.
//  2. Доменные ошибки (PartNotFound, AlreadyPaid, AlreadyCancelled) корректно
//     поднимаются через слои до HTTP-кода.

// TestArch_TotalPrice_ComputedInService проверяет, что цена, которая
// возвращается клиенту, — это в точности сумма цен всех заказанных деталей,
// независимо от их комбинации. Эта инвариант теперь проверяется на
// service-уровне; week_1 считал тот же total в HTTP-хендлере.
func TestArch_TotalPrice_ComputedInService(t *testing.T) {
	cases := []struct {
		name string
		req  *testutil.CreateOrderRequest
		want int64
	}{
		{
			name: "только корпус и двигатель",
			req: &testutil.CreateOrderRequest{
				HullUUID:   testutil.HullAluminumUUID,
				EngineUUID: testutil.EngineIonCUUID,
			},
			want: testutil.HullAluminumPrice + testutil.EngineIonCPrice,
		},
		{
			name: "корпус, двигатель и щит",
			req: &testutil.CreateOrderRequest{
				HullUUID:   testutil.HullTitaniumUUID,
				EngineUUID: testutil.EngineIonBUUID,
				ShieldUUID: testutil.Ptr(testutil.ShieldEnergyUUID),
			},
			want: testutil.HullTitaniumPrice + testutil.EngineIonBPrice + testutil.ShieldEnergyPrice,
		},
		{
			name: "все четыре слота",
			req: &testutil.CreateOrderRequest{
				HullUUID:   testutil.HullTitaniumUUID,
				EngineUUID: testutil.EngineIonBUUID,
				ShieldUUID: testutil.Ptr(testutil.ShieldEnergyUUID),
				WeaponUUID: testutil.Ptr(testutil.WeaponLaserUUID),
			},
			want: testutil.HullTitaniumPrice + testutil.EngineIonBPrice + testutil.ShieldEnergyPrice + testutil.WeaponLaserPrice,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := testutil.NewEnv(t)
			result, resp := env.CreateOrder(t, tc.req)
			defer func() { _ = resp.Body.Close() }()

			require.Equal(t, http.StatusCreated, resp.StatusCode)
			require.NotNil(t, result)
			assert.Equal(t, tc.want, result.TotalPrice,
				"total price считает service-слой как сумму прайсов; HTTP только пробрасывает результат")
		})
	}
}

// Следующие тесты явно демонстрируют, что доменные ошибки разных слоёв
// (client/service/repository) корректно мапятся в HTTP-коды. Это и есть
// главное свойство layered-архитектуры — ошибки не сваливаются в общий 500.

func TestArch_DomainError_HullNotFound_Returns404(t *testing.T) {
	env := testutil.NewEnv(t)
	_, resp := env.CreateOrder(t, &testutil.CreateOrderRequest{
		HullUUID:   uuid.New().String(),
		EngineUUID: testutil.EngineIonCUUID,
	})
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode,
		"PartNotFound из inventory-клиента должен подняться через service до 404")
}

func TestArch_DomainError_PayNonexistentOrder_Returns404(t *testing.T) {
	env := testutil.NewEnv(t)
	_, resp := env.PayOrder(t, uuid.New().String(), &testutil.PayOrderRequest{PaymentMethod: "CARD"})
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode,
		"OrderNotFound из repository должен подняться через service до 404")
}

func TestArch_DomainError_PayPaid_Returns409(t *testing.T) {
	env := testutil.NewEnv(t)
	created, createResp := env.CreateOrder(t, &testutil.CreateOrderRequest{
		HullUUID:   testutil.HullAluminumUUID,
		EngineUUID: testutil.EngineIonCUUID,
	})
	_ = createResp.Body.Close()
	require.NotNil(t, created)

	_, payResp := env.PayOrder(t, created.OrderUUID, &testutil.PayOrderRequest{PaymentMethod: "CARD"})
	_ = payResp.Body.Close()

	_, resp := env.PayOrder(t, created.OrderUUID, &testutil.PayOrderRequest{PaymentMethod: "CARD"})
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusConflict, resp.StatusCode,
		"проверку статуса перед оплатой делает service-слой, не HTTP-хендлер")
}

func TestArch_DomainError_CancelPaid_Returns409(t *testing.T) {
	env := testutil.NewEnv(t)
	created, createResp := env.CreateOrder(t, &testutil.CreateOrderRequest{
		HullUUID:   testutil.HullAluminumUUID,
		EngineUUID: testutil.EngineIonCUUID,
	})
	_ = createResp.Body.Close()
	require.NotNil(t, created)

	_, payResp := env.PayOrder(t, created.OrderUUID, &testutil.PayOrderRequest{PaymentMethod: "CARD"})
	_ = payResp.Body.Close()

	_, resp := env.CancelOrder(t, created.OrderUUID)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusConflict, resp.StatusCode,
		"проверку статуса перед отменой делает service-слой, не HTTP-хендлер")
}

// TestArch_DomainError_OutOfStock_Returns409 проверяет, что ошибка
// ErrOutOfStock из inventory-клиента поднимается через service-слой и
// маппится в 409 (сейчас ErrOutOfStock = Conflict, а не NotFound).
func TestArch_DomainError_OutOfStock_Returns409(t *testing.T) {
	env := testutil.NewEnv(t)
	_, resp := env.CreateOrder(t, &testutil.CreateOrderRequest{
		HullUUID:   testutil.HullOutOfStockUUID,
		EngineUUID: testutil.EngineIonCUUID,
	})
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusConflict, resp.StatusCode,
		"OutOfStock из inventory-клиента должен подняться через service до 409")
}

// TestArch_DomainError_CancelCancelled_Returns409 проверяет, что повторная
// отмена уже отменённого заказа возвращает 409 (ErrOrderCancelled из
// service-слоя), а не 200/500.
func TestArch_DomainError_CancelCancelled_Returns409(t *testing.T) {
	env := testutil.NewEnv(t)
	created, createResp := env.CreateOrder(t, &testutil.CreateOrderRequest{
		HullUUID:   testutil.HullAluminumUUID,
		EngineUUID: testutil.EngineIonCUUID,
	})
	_ = createResp.Body.Close()
	require.NotNil(t, created)

	_, firstResp := env.CancelOrder(t, created.OrderUUID)
	_ = firstResp.Body.Close()

	_, resp := env.CancelOrder(t, created.OrderUUID)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusConflict, resp.StatusCode,
		"повторная отмена уже отменённого заказа должна давать 409 (ErrOrderCancelled)")
}

// TestArch_DomainError_InvalidPaymentMethod_Returns400 проверяет, что
// невалидный enum payment_method отклоняется с 400 (валидация ogen на
// границе API, не доходит до service-слоя).
func TestArch_DomainError_InvalidPaymentMethod_Returns400(t *testing.T) {
	env := testutil.NewEnv(t)
	created, createResp := env.CreateOrder(t, &testutil.CreateOrderRequest{
		HullUUID:   testutil.HullAluminumUUID,
		EngineUUID: testutil.EngineIonCUUID,
	})
	_ = createResp.Body.Close()
	require.NotNil(t, created)

	body := []byte(`{"payment_method": "BITCOIN"}`)
	httpReq, err := http.NewRequest(http.MethodPost,
		env.BaseURL+"/api/v1/orders/"+created.OrderUUID+"/pay",
		bytes.NewReader(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := env.HTTPClient.Do(httpReq)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
		"валидация enum payment_method живёт на границе API (ogen), не в service-слое")
}

// TestArch_DomainError_InvalidUUIDInPath_Returns400 проверяет, что невалидный
// UUID в path отклоняется с 400 (валидация ogen на границе API, не доходит
// до service-слоя и не маппится в 500).
func TestArch_DomainError_InvalidUUIDInPath_Returns400(t *testing.T) {
	env := testutil.NewEnv(t)
	resp, err := env.HTTPClient.Get(env.BaseURL + "/api/v1/orders/not-a-uuid")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
		"валидация формата UUID живёт на границе API (ogen), не в service-слое")
}
