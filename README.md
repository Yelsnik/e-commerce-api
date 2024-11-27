# E-COMMERCE API IN GO 

This API allows users to view products and product details without having to sign up. It allows users with 'merchant' or 'admin' roles to add products, update products or delete products. 

Users who are buyers can add products/items to their cart, update the item and remove it if they so choose.

## Product Endpoints
### Add product
```
POST = {{URL}}/v1/product
```

Remember, for a user to add product, he must be logged in with role set to 'merchant' or 'admin'.

