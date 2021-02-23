# RSA Signer (Prototype)
<h3>Command-line RSA Hash Signer written in Go</h3>

<h5>Usage:</h5>
<pre> -bits int
       Keypair bit length. (for keypair generation only) (default 2048)
 -digest string
       Compute SHA256 hashsum of a file.
 -digest512 string
       Compute SHA512 hashsum of a file.
 -generate
       Generate RSA keypair.
 -hash string
       Input hash/string to sign/verify. (- for stdin)
 -hmac string
       Compute SHA256 HMAC of a file.
 -hmac512 string
       Compute SHA512 HMAC of a file.
 -iter int
       Iterations. (for HMAC only) (default 1)
 -key string
       HMAC secret key.
 -salt string
       Salt. (for HMAC only)
 -sign
       Sign hash with private key.
 -signature string
       Input signature. (verification only)
 -suffix string
       Suffix. (for keypair generation only) (default ".pem")
 -verify
       Verify hash with public key.</pre>
<h5>Example:</h5>
<pre>./rsasigner -digest512 main.go
hash=$(./rsasigner -digest512 main.go)
./rsasigner -sign -key private.pem -hash $hash
sign=$(./rsasigner -sign -key private.pem -hash $hash)
./rsasigner -verify -key public.pem -hash $hash -signature $sign
</pre>
<h5>or:</h5>
<pre>./rsasigner -digest512 main.go|./rsasigner -sign -key private.pem -hash - > sign.txt
sign=$(cat sign.txt)
./rsasigner -digest512 main.go|./rsasigner -verify -key public.pem -hash - -signature $sign
</pre>
<h5>Sign a file:</h5>
<pre>./rsasigner -sign -key private.pem -hash - < file.ext > sign.txt
sign=$(cat sign.txt)
./rsasigner -verify -key public.pem -signature $sign -hash - < file.ext
</pre>
<h5>SHA256-based HMAC with PBKDF2:</h5>
<pre>./rsasigner -hmac &lt;file.ext&gt; -key &lt;secretkey&gt; [-iter 10000|-salt &lt;yoursalt&gt;]
</pre>
</pre>
