# Go-Rest-Api

> Simple RESTful API to create, read, update and delete product, cart, customer, order.

## Quick Start


``` bash
# Install mux router
go get -u github.com/gorilla/mux

#Install gorm
go get -u gorm.io/gorm
```


## Endpoints

### Get All Products
``` bash
GET api/products
```
### Create Cart
``` bash
POST api/carts

# Request sample
# {
#   "id":2,
#   "customerId":2
# }

```

### Get Cart By Id
``` bash
GET api/carts/{id}
```

### Delete Cart Item
``` bash
DELETE api/cartItem/{id}
```

### Create Order
``` bash
POST api/orders/{cartId}
```

## App Info

### Author

Nurşah KARA

### Version

1.0.0
