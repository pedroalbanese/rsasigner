package main    
import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
        b64 "encoding/base64"
)

        var bit = flag.Int("bits", 2048, "Keypair bit length. (for keypair generation only)")
        var digest = flag.String("digest", "", "Compute SHA256 hashsum of a file.")
        var digest512 = flag.String("digest512", "", "Compute SHA512 hashsum of a file.")
        var generate = flag.Bool("generate", false, "Generate RSA keypair.")
        var hash = flag.String("hash", "", "Input hash/string to sign/verify. (- for stdin)")
        var key = flag.String("key", "", "Path to Private/Public key depending on operation.")
        var sig = flag.String("signature", "", "Input signature. (verification only)")
        var sign = flag.Bool("sign", false, "Sign hash with private key.")
        var suf = flag.String("suffix", ".pem", "Suffix. (for keypair generation only)")
        var verify = flag.Bool("verify", false, "Verify hash with public key.")

func main() {
    flag.Parse()

        if (len(os.Args) < 2) {
	fmt.Println("Select -sign, -verify, -generate or -digest.")
        flag.PrintDefaults()
        os.Exit(1)
        } 

        if *sign == false && *verify == false && *generate == false && *digest == "" && *digest512 == "" {
	fmt.Println("Usage:")
	fmt.Println("Select -sign, -verify, -generate or -digest. (type -h)")
        os.Exit(1)
        } 

        if *digest != "" {
	        f, err := os.Open(*digest)
	        if err != nil {
	            log.Fatal(err)
	        }
	            defer f.Close()
	        h := sha256.New()
	        if _, err := io.Copy(h, f); err != nil {
	            log.Fatal(err)
	        }
	        fmt.Printf("%x", h.Sum(nil))
	        os.Exit(0)
        }

        if *digest512 != "" {
	        f, err := os.Open(*digest512)
	        if err != nil {
	            log.Fatal(err)
	        }
	            defer f.Close()
	        h := sha512.New()
	        if _, err := io.Copy(h, f); err != nil {
	            log.Fatal(err)
	        }
	        fmt.Printf("%x", h.Sum(nil))
	        os.Exit(0)
        }

        if *generate == true {
	        GenerateRsaKey(*bit)
	        os.Exit(0)
        }

        if *sign == true && (*key == "" || *hash == "") {
	fmt.Println("Usage:")
	fmt.Println(os.Args[0] + " -sign -key <privatekey.pem> -hash <$hash>")
        os.Exit(1)
        } else if *sign == true && *hash != "-" {
	sourceData := []byte(*hash)
	signData, err := SignatureRSA(sourceData)
	if err != nil {
		 fmt.Println("cryption error:", err)
        os.Exit(1)
	}
        fmt.Print(b64.StdEncoding.EncodeToString(signData))
        os.Exit(0)
	} else if *sign == true && *hash == "-" {
        scannerWrite := bufio.NewScanner(os.Stdin)   		
        if !scannerWrite.Scan() {   			
                log.Printf("Failed to read: %v", scannerWrite.Err()) 
        return
        }
        hash := scannerWrite.Bytes()
	sourceData := []byte(hash)
	signData, err := SignatureRSA(sourceData)
	if err != nil {
		 fmt.Println("cryption error:", err)
        os.Exit(1)
	}
        fmt.Print(b64.StdEncoding.EncodeToString(signData))
        os.Exit(0)
	}

        if *verify == true && (*key == "" || *hash == "" || *sig == "") {
	fmt.Println("Usage:")
	fmt.Println(os.Args[0] + " -verify -key <publickey.pem> -hash <$hash> -signature <$signature>")
        os.Exit(1)
        } else if *verify == true && *hash == "-" {
        scannerWrite := bufio.NewScanner(os.Stdin)   		
        if !scannerWrite.Scan() {   			
                log.Printf("Failed to read: %v", scannerWrite.Err()) 
        return
        }
        hash := scannerWrite.Bytes()
        Signature, err := b64.StdEncoding.DecodeString(*sig)
	err = VerifyRSA([]byte(hash), Signature)
	if err != nil {
		 fmt.Println("Checksum error:", err)
                 os.Exit(0)
        }
	} else if *verify == true && *hash != "-" {
        Signature, err := b64.StdEncoding.DecodeString(*sig)
	err = VerifyRSA([]byte(*hash), Signature)
	if err != nil {
		 fmt.Println("Checksum error:", err)
                 os.Exit(0)
	}
        }
	fmt.Println("Verify correct.")
}

func SignatureRSA(sourceData []byte) ([]byte, error) {
	msg := []byte("")
	file, err := os.Open(*key)
	if err != nil {
		return msg, err
	}
	info, err := file.Stat()
	if err != nil {
		return msg, err
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return msg, err
	}
	myHash := sha256.New()
	myHash.Write(sourceData)
	hashRes := myHash.Sum(nil)
	res, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashRes)
	if err != nil {
		return msg, err
	}
	defer file.Close()
	return res, nil
}

func VerifyRSA(sourceData, signedData []byte) error {
	file, err := os.Open(*key)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	block, _ := pem.Decode(buf)
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	publicKey := publicInterface.(*rsa.PublicKey)
	mySha := sha256.New()
	mySha.Write(sourceData)
	res := mySha.Sum(nil)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, res, signedData)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func GenerateRsaKey(bit int) error {
	private, err := rsa.GenerateKey(rand.Reader, bit)
	if err != nil {
		return err
	}
	privateStream := x509.MarshalPKCS1PrivateKey(private)
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateStream,
	}
	file, err := os.Create("private" + *suf)
	if err != nil {
		return err
	}
	err = pem.Encode(file, &block)
	if err != nil {
		return err
	}
	public := private.PublicKey
	publicStream, err := x509.MarshalPKIXPublicKey(&public)
	if err != nil {
		return err
	}
	pubblock := pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicStream,}
	pubfile, err := os.Create("public" + *suf)
	if err != nil {
		return err
	}
	err = pem.Encode(pubfile, &pubblock)
	if err != nil {
		return err
	}
	return nil
}
