///go:build integration

package httpsrv

import (
	"context"
	"net/http"
	"path"
	"strconv"
	"testing"

	application "github.com/MyLi2tlePony/messenger/internal/app"
	"github.com/MyLi2tlePony/messenger/internal/client"
	databaseConfig "github.com/MyLi2tlePony/messenger/internal/config/database"
	loggerConfig "github.com/MyLi2tlePony/messenger/internal/config/logger"
	serverConfig "github.com/MyLi2tlePony/messenger/internal/config/server"
	"github.com/MyLi2tlePony/messenger/internal/server/http/dto"
	"github.com/MyLi2tlePony/messenger/internal/storage/postgres"
	"github.com/stretchr/testify/require"
)

var configPath = path.Join("..", "..", "..", "configs", "test", "config.toml")

func TestCourier(t *testing.T) {
	dbConfig, err := databaseConfig.New(configPath)
	require.Nil(t, err)

	ctx := context.Background()
	db, err := postgres.Connect(ctx, dbConfig.GetConnectionString())
	require.Nil(t, err)

	apps := application.New(db)

	logConfig, err := loggerConfig.New(configPath)
	require.Nil(t, err)

	srv := NewServer(logConfig.GetLevel(), apps)
	srvConfig, err := serverConfig.New(configPath)
	require.Nil(t, err)

	go func() {
		require.Nil(t, srv.Start(srvConfig.GetHostPort()))
	}()

	httpClient := client.NewHttpClient(http.Client{}, "http://"+srvConfig.GetHostPort())

	userNumber := 6
	usersRequests := make([]dto.CreateUserRequest, userNumber)
	tockens := make([]dto.Tocken, userNumber)
	users := make([]dto.User, userNumber)
	for i := 0; i < userNumber; i++ {
		users[i].Login = "Login" + strconv.Itoa(i)
		users[i].SecondName = "SecondName" + strconv.Itoa(i)
		users[i].FirstName = "FirstName" + strconv.Itoa(i)
		users[i].PublicId = "PublicId" + strconv.Itoa(i)
	}

	for i := 0; i < userNumber; i++ {
		usersRequests[0] = dto.CreateUserRequest{
			Login:    users[i].Login,
			Password: "Password" + strconv.Itoa(i),
		}

		require.Nil(t, httpClient.CreateUser(dto.CreateUserRequest{
			Login:    usersRequests[i].Login,
			Password: usersRequests[i].Password,
		}))

		tockens[i], err = httpClient.CreateTocken(dto.CreateTockenRequest{
			Login:    usersRequests[i].Login,
			Password: usersRequests[i].Password,
		})

		require.Nil(t, err)

		err = httpClient.UpdateUser(dto.UpdateUserRequest{
			Tocken: tockens[i],
			User:   users[i],
		})

		require.Nil(t, err)

		user, err := httpClient.SelectUserByTocken(tockens[i])
		require.Nil(t, err)
		require.Equal(t, user, users[i])

		user, err = httpClient.SelectUserByPublicId(users[i].PublicId)
		require.Nil(t, err)
		require.Equal(t, user, users[i])

		userId, err := db.SelectUserIdByTocken(tockens[i].Text)
		require.Nil(t, err)

		err = db.DeleteUser(userId)
		require.Nil(t, err)
	}
}

//func TestCourier(t *testing.T) {
//	dbConfig, err := databaseConfig.New(configPath)
//	require.Nil(t, err)
//
//	ctx := context.Background()
//	db, err := postgres.Connect(ctx, dbConfig.GetConnectionString())
//	require.Nil(t, err)
//
//	apps := application.New(db)
//
//	logConfig, err := loggerConfig.New(configPath)
//	require.Nil(t, err)
//
//	srv := NewServer(ctx, logConfig.GetLevel(), apps)
//	srvConfig, err := serverConfig.New(configPath)
//	require.Nil(t, err)
//
//	go func() {
//		srv.Start(srvConfig.GetHostPort())
//	}()
//
//	httpClient := &http.Client{}
//	couriersDto := make([]dto.CreateCourierDto, 0)
//	for i := 0; i < 6; i++ {
//		couriersDto = append(couriersDto, dto.CreateCourierDto{
//			CourierType:  "FOOT",
//			Regions:      []int{i, i + 1},
//			WorkingHours: []string{fmt.Sprintf("%d:00-%d:00", i, i+5)},
//		})
//	}
//
//	createCourierRequest := dto.CreateCourierRequest{
//		Couriers: couriersDto,
//	}
//	content, err := json.Marshal(createCourierRequest)
//	require.Nil(t, err)
//
//	url := "http://" + srvConfig.GetHostPort() + urlCouriers
//	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(content))
//	request.Header.Set("Content-Type", "application/json")
//	require.Nil(t, err)
//
//	resp, err := httpClient.Do(request)
//	require.Nil(t, err)
//
//	require.True(t, resp.StatusCode == http.StatusOK)
//	createCouriersResponse := dto.CreateCouriersResponse{}
//
//	content, err = io.ReadAll(resp.Body)
//	require.Nil(t, err)
//
//	require.Nil(t, resp.Body.Close())
//
//	require.Nil(t, json.Unmarshal(content, &createCouriersResponse))
//	couriers := createCouriersResponse.Couriers
//
//	require.True(t, len(couriers) == len(couriersDto))
//	for i := range couriers {
//		require.True(t, couriers[i].EqualsCreateCourierDto(couriersDto[i]))
//
//		id := couriers[i].CourierId
//		url = "http://" + srvConfig.GetHostPort() + urlCouriers + "/" + strconv.Itoa(id)
//		request, err = http.NewRequest(http.MethodGet, url, nil)
//		require.Nil(t, err)
//
//		resp, err = httpClient.Do(request)
//		require.Nil(t, err)
//
//		require.True(t, resp.StatusCode == http.StatusOK)
//		courier := dto.CourierDto{}
//
//		content, err = io.ReadAll(resp.Body)
//		require.Nil(t, err)
//
//		require.Nil(t, resp.Body.Close())
//		require.Nil(t, json.Unmarshal(content, &courier))
//
//		require.True(t, courier.Equals(couriers[i]))
//	}
//
//	url = "http://" + srvConfig.GetHostPort() + urlCouriers
//	request, err = http.NewRequest(http.MethodGet, url, nil)
//	request.Header.Set("Content-Type", "application/json")
//	require.Nil(t, err)
//
//	offset, limit := 2, len(couriers)-2
//	q := request.URL.Query()
//	q.Add("limit", strconv.Itoa(limit))
//	q.Add("offset", strconv.Itoa(offset))
//	request.URL.RawQuery = q.Encode()
//
//	resp, err = httpClient.Do(request)
//	require.Nil(t, err)
//
//	require.True(t, resp.StatusCode == http.StatusOK)
//	getCourierResponse := dto.GetCouriersResponse{}
//
//	content, err = io.ReadAll(resp.Body)
//	require.Nil(t, err)
//
//	require.Nil(t, resp.Body.Close())
//	require.Nil(t, json.Unmarshal(content, &getCourierResponse))
//
//	couriersPart := getCourierResponse.Couriers
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, couriers[i+offset].Equals(couriersPart[i]))
//	}
//
//	for i := range couriers {
//		require.Nil(t, db.DeleteCourierById(ctx, couriers[i].CourierId))
//	}
//
//	for i := range couriers {
//		id := couriers[i].CourierId
//		url = "http://" + srvConfig.GetHostPort() + urlCouriers + "/" + strconv.Itoa(id)
//		request, err = http.NewRequest(http.MethodGet, url, nil)
//		require.Nil(t, err)
//
//		resp, err = httpClient.Do(request)
//		require.Nil(t, err)
//
//		require.True(t, resp.StatusCode == http.StatusNotFound)
//	}
//}
//
//func TestOrder(t *testing.T) {
//	dbConfig, err := databaseConfig.New(configPath)
//	require.Nil(t, err)
//
//	ctx := context.Background()
//	db, err := postgres.Connect(ctx, dbConfig.GetConnectionString())
//	require.Nil(t, err)
//
//	apps := application.New(db)
//
//	logConfig, err := loggerConfig.New(configPath)
//	require.Nil(t, err)
//
//	srv := NewServer(ctx, logConfig.GetLevel(), apps)
//	srvConfig, err := serverConfig.New(configPath)
//	require.Nil(t, err)
//
//	go func() {
//		srv.Start(srvConfig.GetHostPort())
//	}()
//
//	httpClient := &http.Client{}
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
//	createOrderRequest := dto.CreateOrderRequest{
//		Orders: ordersDto,
//	}
//	content, err := json.Marshal(createOrderRequest)
//	require.Nil(t, err)
//
//	url := "http://" + srvConfig.GetHostPort() + urlOrders
//	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(content))
//	request.Header.Set("Content-Type", "application/json")
//	require.Nil(t, err)
//
//	resp, err := httpClient.Do(request)
//	require.Nil(t, err)
//
//	require.True(t, resp.StatusCode == http.StatusOK)
//	orders := make([]dto.OrderDto, 0)
//
//	content, err = io.ReadAll(resp.Body)
//	require.Nil(t, err)
//
//	require.Nil(t, resp.Body.Close())
//
//	require.Nil(t, json.Unmarshal(content, &orders))
//
//	require.True(t, len(orders) == len(ordersDto))
//	for i := range orders {
//		require.True(t, orders[i].EqualsCreateOrderDto(ordersDto[i]))
//
//		id := orders[i].OrderId
//		url = "http://" + srvConfig.GetHostPort() + urlOrders + "/" + strconv.Itoa(id)
//		request, err = http.NewRequest(http.MethodGet, url, nil)
//		require.Nil(t, err)
//
//		resp, err = httpClient.Do(request)
//		require.Nil(t, err)
//
//		require.True(t, resp.StatusCode == http.StatusOK)
//		order := dto.OrderDto{}
//
//		content, err = io.ReadAll(resp.Body)
//		require.Nil(t, err)
//
//		require.Nil(t, resp.Body.Close())
//		require.Nil(t, json.Unmarshal(content, &order))
//
//		require.True(t, order.Equals(orders[i]))
//	}
//
//	url = "http://" + srvConfig.GetHostPort() + urlOrders
//	request, err = http.NewRequest(http.MethodGet, url, nil)
//	request.Header.Set("Content-Type", "application/json")
//	require.Nil(t, err)
//
//	offset, limit := 2, len(orders)-2
//	q := request.URL.Query()
//	q.Add("limit", strconv.Itoa(limit))
//	q.Add("offset", strconv.Itoa(offset))
//	request.URL.RawQuery = q.Encode()
//
//	resp, err = httpClient.Do(request)
//	require.Nil(t, err)
//
//	require.True(t, resp.StatusCode == http.StatusOK)
//
//	content, err = io.ReadAll(resp.Body)
//	require.Nil(t, err)
//
//	ordersPart := make([]dto.OrderDto, 0)
//
//	require.Nil(t, resp.Body.Close())
//	require.Nil(t, json.Unmarshal(content, &ordersPart))
//
//	for i := 0; i+offset < offset+limit; i++ {
//		require.True(t, orders[i+offset].Equals(ordersPart[i]))
//	}
//
//	for i := range orders {
//		require.Nil(t, db.DeleteOrderById(ctx, orders[i].OrderId))
//	}
//
//	for i := range orders {
//		id := orders[i].OrderId
//		url = "http://" + srvConfig.GetHostPort() + urlOrders + "/" + strconv.Itoa(id)
//		request, err = http.NewRequest(http.MethodGet, url, nil)
//		require.Nil(t, err)
//
//		resp, err = httpClient.Do(request)
//		require.Nil(t, err)
//
//		require.True(t, resp.StatusCode == http.StatusNotFound)
//	}
//}

//func TestBalanceTransfer(t *testing.T) {
//	dbConfig, err := databaseConfig.New(configPath)
//	require.Nil(t, err)
//
//	dbConnectionString := dbConfig.GetConnectionString()
//	sqlStorage := balance.Storage(storage.New(dbConnectionString))
//
//	logConfig, err := loggerConfig.New(configPath)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	srvConfig, err := serverConfig.New(configPath)
//	if err != nil {
//		return
//	}
//
//	srv := NewServer(logger.New(logConfig), balance.New(sqlStorage), srvConfig.GetHostPort())
//	go func() {
//		require.Equal(t, http.ErrServerClosed, srv.Start())
//	}()
//
//	defer func() {
//		require.Nil(t, srv.Stop())
//	}()
//
//	httpClient := &http.Client{}
//
//	ctx := context.Background()
//	require.Nil(t, sqlStorage.Connect(ctx))
//
//	srcUserBefore := entity.NewUser(110, 1000)
//	require.Nil(t, sqlStorage.CreateUser(ctx, srcUserBefore))
//
//	dstUserBefore := entity.NewUser(111, 0)
//	require.Nil(t, sqlStorage.CreateUser(ctx, dstUserBefore))
//
//	tb := transferredBalance{
//		SrcUserID:  srcUserBefore.GetID(),
//		DstUserID:  dstUserBefore.GetID(),
//		CreateDate: time.Date(2022, 2, 2, 2, 2, 2, 0, time.UTC),
//		Amount:     1000,
//	}
//	jsonUser, err := json.Marshal(tb)
//	require.Nil(t, err)
//
//	request, err := http.NewRequest(http.MethodPost, "http://"+srvConfig.GetHostPort()+urlBalanceTransfer, bytes.NewBuffer(jsonUser))
//	require.Nil(t, err)
//	request = request.WithContext(ctx)
//
//	resp, err := httpClient.Do(request)
//	require.Nil(t, err)
//
//	require.Nil(t, resp.Body.Close())
//
//	srcUserAfter := entity.NewUser(110, 0)
//	selectedSrcUser, err := sqlStorage.SelectUser(ctx, srcUserAfter.GetID())
//	require.Nil(t, err)
//	require.True(t, reflect.DeepEqual(srcUserAfter, selectedSrcUser))
//
//	dstUserAfter := entity.NewUser(111, 1000)
//	selectedDstUser, err := sqlStorage.SelectUser(ctx, dstUserAfter.GetID())
//	require.Nil(t, err)
//	require.True(t, reflect.DeepEqual(dstUserAfter, selectedDstUser))
//
//	selectedTb, err := sqlStorage.SelectTransferredBalance(ctx, srcUserAfter.GetID())
//	require.Nil(t, err)
//	require.True(t, len(selectedTb) > 0)
//	require.True(t, selectedTb[0].GetSrcUserID() == tb.GetSrcUserID())
//	require.True(t, selectedTb[0].GetAmount() == tb.GetAmount())
//	require.True(t, selectedTb[0].GetDstUserID() == tb.GetDstUserID())
//
//	require.Nil(t, sqlStorage.DeleteUser(ctx, selectedSrcUser.GetID()))
//	require.Nil(t, sqlStorage.DeleteUser(ctx, selectedDstUser.GetID()))
//	require.Nil(t, sqlStorage.DeleteTransferredBalance(ctx, tb.GetSrcUserID(), tb.GetDstUserID(), tb.GetCreateDate()))
//}
