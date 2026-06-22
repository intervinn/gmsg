# api
handles db interaction and http routes

## generating jwt keys:
```bash
# private key (generating)
openssl genrsa -out private.pem 2048

# public key (verifying)
openssl rsa -in private.pem -pubout -out public.pem
```
print one liners to put into env:
```bash
# public key
echo "\"$(awk '{printf "%s\\n", $0}' private.pem)\""

# private key
echo "\"$(awk '{printf "%s\\n", $0}' public.pem)\""
```