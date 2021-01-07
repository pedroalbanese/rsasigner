# rsasigner
<h3>Command-line Hash Signer written in Go</h3>
  
<pre>  
 Usage: 
   -bits int 
        Key pair bit length. (for key pair generation only) (default 2048)
  -digest string 
        Compute SHA256 hashsum of a file.
  -digest512 string <BR>
        Compute SHA512 hashsum of a file. 
  -generate
        Generate RSA key pair. 
  -hash string <BR>
        Input hash/string to sign/verify. (- for stdin) 
  -key string
        Path to Private/Public key depending on operation. 
  -sign 
        Sign hash with private key. 
  -signature string <BR>
        Input signature. (verification only)
  -suffix string <BR>
        Suffix. (for key pair generation only) (default ".pem") 
  -verify <BR>
        Verify hash with public key.
</pre>
