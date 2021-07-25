# Golang-RestAPI
Implementasi berbagai tutorial pembuatan backend restapi dan jwt auth menggunakan go language

User dapat melakukan :
1. Register (http://localhost:9998/register)
2. Login (http://localhost:9998/login)
3. Logout (http://localhost:9998/logout)
4. Akses homepage after login (http://localhost:9998/home)

Server dapat melakukan REST API termasuk:
1. Mendapatkan semua produk (GET) (http://localhost:9999/api/products)
2. Menginputkan data produk baru (POST) (http://localhost:9999/api/products)
3. Mendapatkan produk tertentu (GET) (http://localhost:9999/api/products/{id}) contoh: (http://localhost:9999/api/products/1)
4. Melakukan update data (PUT) (http://localhost:9999/api/products/{id}) contoh: (http://localhost:9999/api/products/2)
5. Melakukan hapus data (DELETE) (http://localhost:9999/api/products{id}) contoh: (http://localhost:9999/api/products/3)
