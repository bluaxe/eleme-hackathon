<!-- 1. mem package to do cache in locally memory -->

<!-- 1. local users to local mem  -->

<!-- 2. food price local mem -->

<!-- 3. token local mem -->

<!-- 4. cart_id local  with userid -->

5. cart foods cache

<!-- 6. order_id local -->

7. food stock spread to servers
	1000 -> 1. 200, 2. 200, 3 200, redic. 400
	mem food 添加获取所有剩余量的接口，修改spreadFood 函数的方法，防止竞态条件
	