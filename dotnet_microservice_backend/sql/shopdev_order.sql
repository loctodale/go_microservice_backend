Create database shopdev_order;
use shopdev_order;

CREATE TABLE orders (
    order_id INT NOT NULL PRIMARY KEY IDENTITY(1,1),
    order_user_id INT NOT NULL,
    order_shipping_id INT NULL,
    order_payment VARCHAR(255),
    order_tracking_number VARCHAR(255),
    order_status VARCHAR(50),
		created_date DATE DEFAULT GETDATE(),
		updated_date DATE DEFAULT GETDATE(),
		deteled_date DATE NULL
);

CREATE TABLE order_detail (
	order_detail_id INT NOT NULL PRIMARY KEY IDENTITY(1,1),
	order_id INT NOT NULL,
	product_id INT NOT NULL,
	quantity INT NOT NULL,
	price_each_item int not null,
	total_price int not null,
	created_date DATE DEFAULT GETDATE(),
	updated_date DATE DEFAULT GETDATE(),
	deteled_date DATE NULL
	FOREIGN KEY (order_id) REFERENCES orders(order_id)
);	