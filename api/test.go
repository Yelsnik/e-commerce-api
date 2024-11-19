package api

/*
user, _ := randomUser(t)

	product := randomProduct(user)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"category":       product.Category,
				"product_name":   product.ProductName,
				"description":    product.Description,
				"brand":          product.Brand,
				"count_in_stock": product.CountInStock,
				"price":          product.Price,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthentication(t, request, tokenMaker, authorizationTypeBearer, user.Email, user.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateProductsParams{
					ID:           util.Test(),
					Category:     product.Category,
					ProductName:  product.ProductName,
					Description:  product.Description,
					Brand:        product.Brand,
					CountInStock: product.CountInStock,
					Price:        product.Price,
					Rating:       product.Rating,
					UserID:       product.UserID,
				}
				store.EXPECT().CreateProducts(gomock.Any(), gomock.Eq(args)).Times(1).Return(product, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				fmt.Println(recoder.Body)
				requireBodyMatchProduct(t, recoder.Body, product)
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

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/v1/product"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
			fmt.Println(data, request, recorder.Body)
		})
	}
*/
