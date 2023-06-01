//go:build integration

package app

import (
	"path"
)

var configPath = path.Join("..", "..", "configs", "test", "config.toml")

//func TestCouriers(t *testing.T) {
//	dbConfig, err := databaseConfig.New(configPath)
//	require.Nil(t, err)
//
//	ctx := context.Background()
//	db, err := postgres.Connect(ctx, dbConfig.GetConnectionString())
//	require.Nil(t, err)
//
//	apps := New(db)
//
//	couriersDto := make([]dto.CreateCourierDto, 0)
//	for i := 0; i < 6; i++ {
//		couriersDto = append(couriersDto, dto.CreateCourierDto{
//			CourierType:  "FOOT",
//			Regions:      []int{i, i + 1},
//			WorkingHours: []string{fmt.Sprintf("%d:00-%d:00", i, i+5)},
//		})
//	}
//
//	couriers, err := apps.AddCouriers(ctx, couriersDto)
//	require.Nil(t, err)
//	require.True(t, len(couriers) == len(couriersDto))
//
//	for i := range couriers {
//		require.True(t, couriers[i].EqualsCreateCourierDto(couriersDto[i]))
//
//		courier, err := apps.GetCourierById(ctx, couriers[i].CourierId)
//		require.Nil(t, err)
//		require.True(t, couriers[i].Equals(courier))
//	}
//
//	offset, limit := 0, 2
//	couriersPart, err := apps.GetCouriers(ctx, offset, limit)
//	require.Nil(t, err)
//
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, couriers[i+offset].Equals(couriersPart[i]))
//	}
//
//	offset, limit = 2, len(couriers)-2
//	couriersPart, err = apps.GetCouriers(ctx, offset, limit)
//	require.Nil(t, err)
//
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, couriers[i+offset].Equals(couriersPart[i]))
//	}
//
//	for i := range couriers {
//		require.Nil(t, db.DeleteCourierById(ctx, couriers[i].CourierId))
//	}
//}
//
//func TestOrders(t *testing.T) {
//	dbConfig, err := databaseConfig.New(configPath)
//	require.Nil(t, err)
//
//	ctx := context.Background()
//	db, err := postgres.Connect(ctx, dbConfig.GetConnectionString())
//	require.Nil(t, err)
//
//	apps := New(db)
//
//	ordersDto := make([]dto.CreateOrderDto, 0)
//	for i := 0; i < 6; i++ {
//		ordersDto = append(ordersDto, dto.CreateOrderDto{
//			Weight:        0.25 + float64(i),
//			Regions:       10 + i,
//			DeliveryHours: []string{fmt.Sprintf("%d:00-%d:00", i, i+5)},
//			Cost:          i * 100,
//		})
//	}
//
//	orders, err := apps.AddOrders(ctx, ordersDto)
//	require.Nil(t, err)
//	require.True(t, len(orders) == len(ordersDto))
//
//	for i := range orders {
//		require.True(t, orders[i].EqualsCreateOrderDto(ordersDto[i]))
//
//		courier, err := apps.GetOrderById(ctx, orders[i].OrderId)
//		require.Nil(t, err)
//		require.True(t, orders[i].Equals(courier))
//	}
//
//	offset, limit := 0, 2
//	couriersPart, err := apps.GetOrders(ctx, offset, limit)
//	require.Nil(t, err)
//
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, orders[i+offset].Equals(couriersPart[i]))
//	}
//
//	offset, limit = 2, len(orders)-2
//	couriersPart, err = apps.GetOrders(ctx, offset, limit)
//	require.Nil(t, err)
//
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, orders[i+offset].Equals(couriersPart[i]))
//	}
//
//	for i := range orders {
//		require.Nil(t, db.DeleteOrderById(ctx, orders[i].OrderId))
//	}
//}
