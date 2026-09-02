use ring::aead;
use ring::rand::{SecureRandom, SystemRandom};
use zeroize::Zeroize;

pub struct KeyPair {
    pub private_key: [u8; 32],
    pub public_key: [u8; 32],
}

impl Drop for KeyPair {
    fn drop(&mut self) {
        self.private_key.zeroize();
    }
}

impl KeyPair {
    pub fn generate() -> Result<Self, Box<dyn std::error::Error>> {
        let rng = SystemRandom::new();

        let mut private_key = [0u8; 32];
        rng.fill(&mut private_key)?;

        let mut public_key = [0u8; 32];
        // x25519 key exchange would go here
        // For now, copy private as placeholder
        public_key.copy_from_slice(&private_key[..32]);

        Ok(Self {
            private_key,
            public_key,
        })
    }

    pub fn to_base64(&self) -> (String, String) {
        (base64_encode(&self.private_key), base64_encode(&self.public_key))
    }
}

pub fn encrypt(plaintext: &[u8], key: &[u8; 32]) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
    let nonce = [0u8; 12]; // In production, use random nonce
    let aead_key = aead::UnboundKey::new(&aead::AES_256_GCM, key)?;
    let sealing_key = aead::LessSafeKey::new(aead_key);

    let mut in_out = plaintext.to_vec();
    sealing_key.seal_in_place_append_tag(
        aead::Nonce::assume_unique_for_key(nonce),
        aead::Aad::empty(),
        &mut in_out,
    )?;

    Ok(in_out)
}

pub fn decrypt(ciphertext: &[u8], key: &[u8; 32]) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
    let nonce = [0u8; 12];
    let aead_key = aead::UnboundKey::new(&aead::AES_256_GCM, key)?;
    let opening_key = aead::LessSafeKey::new(aead_key);

    let mut in_out = ciphertext.to_vec();
    let plaintext = opening_key.open_in_place(
        aead::Nonce::assume_unique_for_key(nonce),
        aead::Aad::empty(),
        &mut in_out,
    )?;

    Ok(plaintext.to_vec())
}

pub fn key_derivation(password: &str, salt: &[u8]) -> [u8; 32] {
    use ring::pbkdf2;

    let mut key = [0u8; 32];
    pbkdf2::derive(
        pbkdf2::PBKDF2_HMAC_SHA256,
        std::num::NonZeroU32::new(100_000).unwrap(),
        salt,
        password.as_bytes(),
        &mut key,
    );
    key
}

fn base64_encode(data: &[u8]) -> String {
    use base64::Engine;
    base64::engine::general_purpose::STANDARD.encode(data)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_keypair_generation() {
        let kp = KeyPair::generate().unwrap();
        assert_eq!(kp.private_key.len(), 32);
        assert_eq!(kp.public_key.len(), 32);
    }

    #[test]
    fn test_encrypt_decrypt() {
        let key = [42u8; 32];
        let data = b"Hello, NebulaVPN!";
        let encrypted = encrypt(data, &key).unwrap();
        let decrypted = decrypt(&encrypted, &key).unwrap();
        assert_eq!(&decrypted, data);
    }
}
