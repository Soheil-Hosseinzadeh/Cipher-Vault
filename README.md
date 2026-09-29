# Cipher Vault

Cipher Vault is a desktop password vault built with Go and Fyne. It encrypts vault data with **AES-256 in Galois/Counter Mode (GCM)**. A 32-byte key is derived from the user's passphrase using **Argon2id**, with a fresh random 16-byte salt for each encryption. A random GCM nonce is also generated; the salt and nonce are stored with the ciphertext. GCM authenticates the encrypted data so tampering is detected during decryption.

I developed this program offline and later lost the final source code. The code available here is an imperfect, incomplete reconstruction of the Fyne GUI. The encryption implementation is intact. I was able to recover the compiled Linux and Windows executables, which are included in the `bin` directory.

#### How to work with it?
It's simple.
- Make a new Vault with a strong passkey.
- Store your passwords in it.
- Connect to your Vault at any other time in order to view or update it.

#### Home not connected:
<img width="1150" height="687" alt="image" src="https://github.com/user-attachments/assets/b5faae50-4d83-4cb1-902d-33e73a7ec5f5" />


#### New Vault:
<img width="1150" height="687" alt="image" src="https://github.com/user-attachments/assets/1d7292d7-48ea-4c9d-b1c6-e6609614fa28" />

#### Home connected to a vault:
<img width="1150" height="687" alt="image" src="https://github.com/user-attachments/assets/5c46b46b-799f-45c1-9e71-cf6aa79ce42b" />

#### Connect to a vault:
<img width="1150" height="687" alt="image" src="https://github.com/user-attachments/assets/488d7047-b3c6-414b-b489-a3a9a080c972" />
