package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"project/handler"
	"project/ikan"
	"project/user"
	"strconv"
	"time"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID             uint   `gorm:"primarykey" json:"id"`
	Name           string `json:"nama"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Alamat         string `json:"alamat"`
	Jenis_Budidaya string `json:"jenis_budidaya"`
	Lokasi_Tambak  string `json:"lokasi_tambak"`
	Luas_Kolam     string `json:"luas_kolam"`
	TokoID         string `json:"toko_id"`
}

type Ikan struct {
	ID               uint   `gorm:"primarykey" json:"id"`
	Kategori         string `json:"kategori"`
	Jenis_Ikan       string `json:"jenis_ikan"`
	Harga            string `json:"harga"`
	TokoID           uint   `json:"toko_id"`
	Provinsi         string `json:"provinsi"`
	Kota             string `json:"kota"`
	Bulan_Panen      string `json:"bulan_panen"`
	Deskripsi_Produk string `json:"deskripsi_produk"`
	Toko Toko
}

type Pakan struct {
	ID               uint   `gorm:"primarykey" json:"id"`
	Berat uint `json:"berat"`
	Kategori string `json:"kategori"`
	Etalase string `json:"etalase"`
	Deskripsi string `json:"deskripsi"`
	Kemasan string `json:"kemasan"`
	Bahan_Bahan string `json:"bahan_bahan"`
	Komposisi string `json:"komposisi"`
}


type Toko struct {
	ID uint `gorm:"primarykey" json:"id"`
	NamaToko string `json:"nama_toko"`
}

type Tweet struct {
	ID        uint `gorm:"primarykey" json:"id"`
	UserID    uint `json:"user_id"`
	User      User
	Content   string        `json:"name"`
	RepliedTo sql.NullInt64 `json:"replied_to"`
	CreatedAt time.Time     `json:"created_at"`
}

var db *gorm.DB
var r *gin.Engine

func InitDB() error {
	_db, err := gorm.Open(mysql.Open("root:spenesa234@tcp(127.0.0.1:3306)/intern_workshop?parseTime=true"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	err = db.AutoMigrate(&User{}, &Tweet{}, &Ikan{})

	tweet := Tweet{
		ID: 1,
	}
	db.Preload("User").Take(&tweet)
	fmt.Println(tweet)
	if err != nil {
		return err
	}
	return nil
}

func InitGin() {
	r = gin.Default()
	r.Use(CORSPreflightMiddleware())
}

func CORSPreflightMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	}
}

type postLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type postRegisterToko struct {
	ID uint
}

type patchUserBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		header = header[len("Bearer "):]
		token, err := jwt.Parse(header, func(t *jwt.Token) (interface{}, error) {
			return []byte("passwordBuatSigning"), nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "JWT validation error.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("id", claims["id"])
			c.Next()
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "JWT invalid.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
	}
}

func InitRouter() {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	ikanRepository := ikan.NewRepository(db)
	ikanService := ikan.NewService(ikanRepository)
	ikanHandler := handler.NewIkanHandler(ikanService)

	r.POST("/api/auth/register", userHandler.Register)
	r.POST("/api/auth/tambah-ikan", ikanHandler.RegisterIkan)

	r.POST("/api/auth/register-member", func(c *gin.Context) {
		var body User
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := User{
			Alamat:         body.Alamat,
			Jenis_Budidaya: body.Jenis_Budidaya,
			Lokasi_Tambak:  body.Lokasi_Tambak,
			Luas_Kolam:     body.Luas_Kolam,
		}
		if result := db.Create(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Berhasil Membuat Akun Member",
			"status":  "Sukses",
			"data": gin.H{
				"id":             user.Alamat,
				"Jenis Budidaya": user.Jenis_Budidaya,
				"Lokasi Tambak":  user.Lokasi_Tambak,
				"Luas Kolam":     user.Luas_Kolam,
			},
		})
	})

	r.POST("/api/auth/login", func(c *gin.Context) {
		var body postLoginBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := User{}
		if result := db.Where("email = ?", body.Email).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if user.Password == body.Password {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id":  user.ID,
				"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
			})
			tokenString, err := token.SignedString([]byte("passwordBuatSigning"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when generating the token.",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Login Berhasil",
				"status":  "Sukses",
				"data": gin.H{
					"id":    user.ID,
					"name":  user.Name,
					"token": tokenString,
				},
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Password is incorrect.",
			})
			return
		}
	})

	r.PATCH("/api/auth/update-user", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		var body patchUserBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := User{
			ID:       uint(parsedId),
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
		}
		result := db.Model(&user).Updates(user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result = db.Where("id = ?", parsedId).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Berhasil Memperbarui Akun",
			"status":  "Sukses",
			"data":    user,
		})
	})

	r.GET("/api/auth/ikan-segar", func(c *gin.Context) {
		// ikan := Ikan{}
		var ikan []Ikan
		if result := db.Where("kategori = ?", "ikan segar").Find(&ikan); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Pencarian Berhasil",
			"status":  "Sukses",
			"data":    ikan,
		})
	})

	r.GET("/api/auth/ikan-frozen", func(c *gin.Context) {
		ikan := Ikan{}
		if result := db.Where("kategori = ?", "ikan frozen").Find(&ikan); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Pencarian Berhasil",
			"status":  "Sukses",
			"data":    ikan,
		})
	})

	r.GET("/api/auth/bibit-ikan", func(c *gin.Context) {
		ikan := Ikan{}
		if result := db.Where("kategori = ?", "bibit ikan").Find(&ikan); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Pencarian Berhasil",
			"status":  "Sukses",
			"data":    ikan,
		})
	})

	r.DELETE("/api/auth/delete-ikan", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}
		ikan := Ikan{
			ID: uint(parsedId),
		}
		if result := db.Delete(&ikan); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Hapus Ikan Berhasil",
			"status":  "Sukses",
		})
	})
}

func StartServer() error {
	return r.Run(":5000")
}

func main() {
	if err := InitDB(); err != nil {
		fmt.Println("Database error on init!")
		fmt.Println(err.Error())
		return
	}

	InitGin()
	InitRouter()
	CORSPreflightMiddleware()

	if err := StartServer(); err != nil {
		fmt.Println("Server error!")
		fmt.Println(err.Error())
		return
	}
}
