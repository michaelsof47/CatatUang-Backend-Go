# CatatUang-Backend-Go

CatatUang Backend merupakan migrasi dari backend node.js ke bahasa go dimana concurrency merupakan pertimbangan bagi aku untuk migrasi tersebut. Berikut ini merupakan panduan penggunaan api nya :

## Panduan API Route

1. **/create_account** : digunakan untuk akses buat akun
2. **/login_user** : digunakan untuk akses login akun dan ada beberapa fitur yang diterapkan :
   - Jika User login lebih dari 5 kali. Maka, User tidak dapat spam login dan harus menunggu 5 menit kemudian
   - Jika User login lebih dari 1 device. Maka, User tidak dapat login di device tersebut sehingga harus logout dari device sebelumnya
3. **/get_user/:id** : digunakan untuk mendapatkan data akun berdasarkan id
4. **/logout_user** : digunakan untuk logout akun user termasuk hapus data token dari redis

## Note

Untuk CatatUang Backend masih dalam tahap migrasi dan integrasi sehingga masih dikembangkan sampai saat ini. Saran dan Masukkan sangat membantu saya dalam mengembangkan aplikasi ini.