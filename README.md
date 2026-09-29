# Cipher Vault

Cipher Vault is a desktop password vault built with Go and Fyne. It encrypts vault data with **AES-256 in Galois/Counter Mode (GCM)**. A 32-byte key is derived from the user's passphrase using **Argon2id**, with a fresh random 16-byte salt for each encryption. A random GCM nonce is also generated; the salt and nonce are stored with the ciphertext. GCM authenticates the encrypted data so tampering is detected during decryption.

I developed this program offline and later lost the final source code. The code available here is an imperfect, incomplete reconstruction of the Fyne GUI. The encryption implementation is intact. I was able to recover the compiled Linux and Windows executables, which are included in the `bin` directory.
