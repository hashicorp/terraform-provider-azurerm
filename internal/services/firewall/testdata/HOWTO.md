# How Certificates were generated

## How Key and Certificate was generated
```bash
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem
```

## How PFX was generated from above Key and Certificate
```bash
openssl pkcs12 -export -out certificate.pfx -inkey key.pem -in cert.pem
```