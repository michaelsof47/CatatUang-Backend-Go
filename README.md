# CatatUang-Backend-Go

CatatUang Backend merupakan migrasi dari backend node.js ke bahasa go dimana concurrency merupakan pertimbangan bagi aku untuk migrasi tersebut. Berikut ini merupakan panduan penggunaan api nya :

## Panduan API Route


### Koleksi User

1. `/create_account` : digunakan untuk akses buat akun
   - Apabila ketika yang diinput adalah nama depan dan nama belakang yang sama dengan yang diinput sebelumnya. Maka, akan muncul pesan error kalo harus memasukkan data baru
2. `/login_user` : digunakan untuk akses login akun dan ada beberapa fitur yang diterapkan :
   - Jika User login lebih dari 5 kali. Maka, User tidak dapat spam login dan harus menunggu 5 menit kemudian
   - Jika User login lebih dari 1 device. Maka, User tidak dapat login di device tersebut sehingga harus logout dari device sebelumnya
3. `/get_user/:id` : digunakan untuk mendapatkan data akun berdasarkan id
4. `/logout_user` : digunakan untuk logout akun user termasuk hapus data token dari redis
5. `/update_password` : digunakan untuk update password berdasarkan id
6. `/update_user` : digunakan untuk update informasi user
7. `/update_profile_image` : digunakan untuk update informasi foto profile
8. `/verify_otp` : digunakan untuk verifikasi otp setelah dilakukan login akun
9. `/profile_image` : digunakan untuk mengambil image profile user
10 `/forgot_password` : digunakan untuk update password jika lupa password

### Koleksi Balance

1. `/create_or_update_amount` : digunakan untuk pengisian awal saldo transaksi / update saldo (khusus di api)
2. `/add_balance_amount` : digunakan untuk tambah saldo berdasarkan user_id jika saldo kurang
3. `/get_balance` : digunakan untuk mendapatkan saldo berdasarkan user_id

### Koleksi Transaksi
1. `/create_category` : digunakan untuk menambah kategori berdasarkan user_id
   - Apabila ketika yang diinput adalah nama kategori yang sama dengan yang diinput sebelumnya. Maka, akan muncul pesan error kalo harus memasukkan data baru
2. `/get_categories` : digunakan untuk menampilkan isi kategori dalam bentuk list berdasarkan user_id

## Note

Untuk CatatUang Backend masih dalam tahap migrasi dan integrasi sehingga masih dikembangkan sampai saat ini. Saran dan Masukkan sangat membantu saya dalam mengembangkan aplikasi ini.