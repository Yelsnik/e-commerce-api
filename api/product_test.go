package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Yelsnik/e-commerce-api/db/mock"
	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	//	"github.com/Yelsnik/e-commerce-api/token"
	"github.com/Yelsnik/e-commerce-api/util"
	//"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomProduct(user db.User) db.Product {
	return db.Product{
		//ID:           util.Test(),
		Category:     util.RandomCategory(),
		ProductName:  util.RandomString(6),
		Description:  util.RandomString(7),
		Brand:        util.NewNullString(util.RandomString(6)),
		CountInStock: util.RandomInt(0, 100),
		Price:        float64(util.RandomMoney()),
		Currency:     util.RandomCurrency(),
		Rating:       util.NewNullInt(1),
		IsFeatured:   util.NewNullBool(true),
		UserID:       user.ID,
		CreatedAt:    time.Now(),
	}
}

type response struct {
	Data db.Product `json:"data"`
}

func checkCreateProductResponse(t *testing.T, body []byte, product db.Product) {
	var bites response

	err := json.Unmarshal(body, &bites)
	require.NoError(t, err)
	require.NotEmpty(t, bites)

	require.Equal(t, product.ID, bites.Data.ID)
	require.Equal(t, product.Category, bites.Data.Category)
	require.Equal(t, product.ProductName, bites.Data.ProductName)
	require.Equal(t, product.Description, bites.Data.Description)
	require.Equal(t, product.Brand, bites.Data.Brand)
	require.Equal(t, product.CountInStock, bites.Data.CountInStock)
	require.Equal(t, product.Price, bites.Data.Price)
	require.Equal(t, product.Currency, bites.Data.Currency)
	require.Equal(t, product.Rating, bites.Data.Rating)
	require.Equal(t, product.IsFeatured, bites.Data.IsFeatured)
	require.Equal(t, product.UserID, bites.Data.UserID)
	require.NotZero(t, bites.Data.CreatedAt)
}

func requireBodyMatchProduct(t *testing.T, body *bytes.Buffer, product db.Product) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotProduct db.Product
	err = json.Unmarshal(data, &gotProduct)
	require.NoError(t, err)
	require.Equal(t, product.Category, gotProduct.Category)
}

func requireBodyMatchProducts(t *testing.T, body *bytes.Buffer, products []db.Product) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	fmt.Println(products)

	var gotProducts []db.Product
	err = json.Unmarshal(data, &gotProducts)
	require.NoError(t, err)
	require.Equal(t, products, gotProducts)
}

func TestGetProductApi(t *testing.T) {
	user, _ := randomUser(t)
	product := randomProduct(user)

	testCases := []struct {
		name          string
		productID     uuid.UUID
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			productID: product.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetProducts(gomock.Any(), gomock.Eq(product.ID)).Times(1).Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchProduct(t, recorder.Body, product)
			},
		},
		{
			name:      "invalid id",
			productID: product.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetProducts(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// create a newcontroller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create new store and  build your stubs
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start a new server and a new recorder
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/v1/product/%s", tc.productID.String())
			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, http.StatusOK, recorder.Code)
		})
	}

}

func TestCreateProductApi(t *testing.T) {
	user, _ := randomUser(t)
	product := randomProduct(user)

	// create a newcontroller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// create new store and  build your stubs
	store := mockdb.NewMockStore(ctrl)

	arg := db.CreateProductsParams{
		Category:     product.Category,
		ProductName:  product.ProductName,
		Description:  product.Description,
		Brand:        product.Brand,
		CountInStock: product.CountInStock,
		Price:        product.Price,
		Currency:     product.Currency,
		Rating:       product.Rating,
		IsFeatured:   product.IsFeatured,
		UserID:       product.UserID,
	}

	store.EXPECT().CreateProducts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(product, nil)

	server := newTestServer(t, store)

	recorder := httptest.NewRecorder()

	var str string
	if product.Brand.Valid {
		str = product.Brand.String
	}

	body := gin.H{
		"category":       product.Category,
		"product_name":   product.ProductName,
		"description":    product.Description,
		"brand":          str,
		"count_in_stock": product.CountInStock,
		"price":          product.Price,
		"currency":       product.Currency,
		"rating":         1,
		"is_featured":    true,
	}

	data, err := json.Marshal(body)
	require.NoError(t, err)

	url := "/v1/product"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	addAuthentication(t, request, server.tokenMaker, authorizationTypeBearer, user.ID, user.Role, time.Minute)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.NotEmpty(t, recorder.Body)
	require.NotEmpty(t, recorder.Body.Bytes())

	fmt.Println("1", recorder.Body)
	checkCreateProductResponse(t, recorder.Body.Bytes(), product)

	//requireBodyMatchProduct(t, recorder.Body, product)
}

func TestListProductApi(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	products := make([]db.Product, n)
	for i := 0; i < n; i++ {
		products[i] = randomProduct(user)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListProductsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().ListProducts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(products, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/v1/products"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
