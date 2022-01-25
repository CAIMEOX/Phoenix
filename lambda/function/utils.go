package function

import (
"crypto/sha256"
"encoding/hex"
	"github.com/pelletier/go-toml"
	"io"
	"io/ioutil"
	"log"
	"os"
"strconv"
)

func SliceAtoi(sa []string) ([]float64, error) {
	si := make([]float64, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}

func GetHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func AddVector(v1, v2 Vector) Vector {
	v1[0] += v2[0]
	v1[1] += v2[1]
	v1[2] += v2[2]
	return v1
}

type config struct {
	Connection struct {
		RemoteAddress string
	}
	User struct {
		Bot string
		Auth bool
		Operator string
	}
	Debug struct{
		Enabled bool
	}
	Lib struct {
		Std bool
		Script []string
	}
}

func ReadConfig(path string) config {
	c := config{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("error creating config: %v", err)
		}
		data, err := toml.Marshal(c)
		if err != nil {
			log.Fatalf("error encoding default config: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("error writing encoded default config: %v", err)
		}
		_ = f.Close()
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		log.Fatalf("error decoding config: %v", err)
	}
	data, _ = toml.Marshal(c)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		log.Fatalf("error writing config file: %v", err)
	}
	return c
}
