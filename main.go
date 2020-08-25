package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/outakujo/utils"
	"gopkg.in/gomail.v2"
	"net/http"
	"strconv"
	"time"
)
type bulletScreen struct{
	Id  int  `db:"id"`
	QQ  string  `db:"qq"`
	UserName  string  `db:"userName"`
	BulletScreenText  string  `db:"bulletScreenText"`
}
type TalkStruct struct{
	Id  int  `db:"id"`
	Name  string  `db:"name"`
	TimeSite  string  `db:"timeSite"`
	ViewsNumber  int  `db:"viewsNumber"`
	Comment  int  `db:"comment"`
	Like  int  `db:"likeNum"`
	Text  string  `db:"text"`
}
type TalkComment struct{
	CommentId int `db:"commentId"`
	TalkId int `db:"talkId"`
	ReplyId int `db:"replyId"`
	PositionFlag string `db:"positionFlag"`
	CommentName string `db:"commentName"`
	ObjName string `db:"objName"`
	QQ string `db:"QQ"`
	CommentText string `db:"commentText"`
	CommentTime string `db:"commentTime"`
}
type userInformationStruct struct{
	Id int
	Email string
	SecurityCode int
	Username string
	Password string
	QQ string
}
func verify(c *gin.Context){
	username := c.PostForm("username")
	password := c.PostForm("password")  // 取到就返回值，取不到返回空字符串
	db := mysqlInit()
	defer db.Close()
	var num  int
	err := db.QueryRow("select count(*) from user_information where (username=? or email=?) and password=?", username,username,password).Scan(&num)
	fmt.Println(num)
	if err != nil{
		return
	}
	fmt.Println("NUM",num)
	if num > 0 {
		c.Next()
	} else{
		c.Abort()
	}
}
func main(){
	r := gin.Default()
	r.Static("/css","./css")
	r.Static("/img","./img")
	r.Static("/live2d", "./live2d")
	r.Static("/music", "./music")
	r.Static("/fonts", "./fonts")
	r.LoadHTMLFiles("./login.html", "./index.html","./live2d.html","./算法.html","./bulletScreen.html","./深度学习.html")
	r.GET("",func(c *gin.Context){
		fmt.Printf("2323323")
		c.HTML(http.StatusOK,"login.html",nil)
	})
	r.POST("/login",verify, func(c *gin.Context) {
		// 获取form表单提交的数据
		/*
		typenum := c.PostForm("type")
		emailsign := c.PostForm("emailsign")
		email := c.PostForm("email")
		securitycode := c.PostForm("securitycode")
		username := c.PostForm("username")
		password := c.PostForm("password")  // 取到就返回值，取不到返回空字符串
		qq := c.PostForm("qq")

		//c.String(http.StatusOK, "wangwenjie")
		fmt.Println("typenum:"+typenum)
		fmt.Println("emailsign"+emailsign)
		fmt.Println("email:"+ email)
		fmt.Println("SecurityCode:"+securitycode)
		fmt.Println("username:"+ username)
		fmt.Println("password:"+ password)
		fmt.Println("qq:"+ qq)
		//username := c.DefaultPostForm("username", "somebody")
		//password := c.DefaultPostForm("xxx", "***")
		 */
		username := c.PostForm("username")
		db := mysqlInitx()
		defer db.Close()
		var userInformation []userInformationStruct
		err := db.Select(&userInformation,"select * from user_information where (username=? or email=?)",username,username)
		if err != nil{
			fmt.Println("exec faided",err)
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"username": userInformation[0].Username,
			"password": userInformation[0].Password,
			"qq":userInformation[0].QQ,
		})

		t1:=time.Now().Year()        //年
		t2:=time.Now().Month()       //月
		t3:=time.Now().Day()         //日
		t4:=time.Now().Hour()        //小时
		t5:=time.Now().Minute()      //分钟
		t6:=time.Now().Second()      //秒
		t7:=time.Now().Nanosecond()  //纳秒
		loginTime:=time.Date(t1,t2,t3,t4,t5,t6,t7,time.Local)
		r, err := db.Exec("insert into loginTable(userName,qq,loginTime)values( ?, ?, ?)", userInformation[0].Username,userInformation[0].QQ, loginTime)

		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		id, err := r.LastInsertId()
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		fmt.Println("insert succ:", id)
		fmt.Println(userInformation)

	})
	r.GET("/login1",func(c *gin.Context){
		var passwordsql string
		username := c.Query("username")
		password := c.Query("password")
		db := mysqlInit()
		defer db.Close()
		row := db.QueryRow("select password from user_information where username=? or email=?", username,username)
		err := row.Scan(&passwordsql)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "exist"})
		} else if passwordsql == password{
			c.JSON(http.StatusOK, gin.H{"message": "yes"})
		}else if(passwordsql != password){
			c.JSON(http.StatusOK, gin.H{"message": "no"})
		}
	})
	r.GET("/login2",func(c *gin.Context){
		var emailnum string
		var usernamenum string
		email := c.Query("email")
		securitycode := c.Query("securitycode")
		username := c.Query("username")
		password := c.Query("password")
		qq := c.Query("qq")
		fmt.Println(username)
		db := mysqlInit()
		defer db.Close()
		db.QueryRow("select count(*) from user_information where email=?", email).Scan(&emailnum)
		db.QueryRow("select count(*) from user_information where username=?", username).Scan(&usernamenum)
		n ,err := strconv.Atoi(emailnum)
		if err !=nil{
			return
		}
		m ,err := strconv.Atoi(usernamenum)
		if err !=nil{
			return
		}
		if n > 0{
			c.JSON(http.StatusOK, gin.H{"message": "emailexist"})
		} else if m > 0{
			c.JSON(http.StatusOK, gin.H{"message": "usernameexist"})
		}else{
			t3:=time.Now().Day()
			t4:=time.Now().Hour()
			if securitycode == strconv.Itoa(t3 + t4 + 10000) {
				stmt, err := db.Prepare("INSERT INTO user_information SET email=?,securitycode=?,username=?,password=?,qq=?")
				if err != nil {
					fmt.Println(err)
					return
				}
				res, err := stmt.Exec(email,securitycode,username, password, qq)
				id, err := res.LastInsertId()
				if err != nil {
					panic(err)
				}
				fmt.Println(id)
				c.JSON(http.StatusOK, gin.H{"message": "yes","id":id})
			}else{
				c.JSON(http.StatusOK, gin.H{"message": "nosecuritycode"})
			}
		}
	})
	r.GET("/login3",func(c *gin.Context){
		var securitycodesql string
		email := c.Query("email")
		securitycode := c.Query("securitycode")
		password := c.Query("password")
		db := mysqlInit()
		defer db.Close()
		err := db.QueryRow("select securitycode from user_information where email=?", email).Scan(&securitycodesql)
		if err != nil{
			c.JSON(http.StatusOK, gin.H{"message": "noemail"})
		}
		if securitycodesql != securitycode && err == nil{
			c.JSON(http.StatusOK, gin.H{"message": "nosecuritycode"})
		}
		if securitycodesql == securitycode{
			stmt, err := db.Prepare("update user_information SET password=? where email=?")
			if err != nil {
				fmt.Println(err)
				return
			}
			res, err := stmt.Exec(password,email)
			id, err := res.RowsAffected()
			if err != nil {
				panic(err)
			}
			fmt.Println(id)
			c.JSON(http.StatusOK, gin.H{"message": "yes"})
		}
	})
	r.GET("/loginEmail", func(c *gin.Context) {
		var securitycode int
		db := mysqlInit()
		defer db.Close()
		email := c.Query("email")
		typenum := c.Query("typenum")
		row := db.QueryRow("select securitycode from user_information where email=?", email)
		err := row.Scan(&securitycode)
		if err!=nil {
			if(typenum == "2"){
				t3:=time.Now().Day()
				t4:=time.Now().Hour()
				sendEmail(email,t3+t4+10000)
				db.QueryRow("select securitycode from user_information where email=?", email)
			}
			c.JSON(http.StatusOK, gin.H{"message": "no"})
		}else{
			if(typenum == "3"){
				sendEmail(email,securitycode)
			}
			c.JSON(http.StatusOK, gin.H{"message": "yes"})
		}

	})
	r.GET("/talk1",func(c *gin.Context){
		db := mysqlInitx()
		defer db.Close()
		var talk []TalkStruct
		err := db.Select(&talk,"select * from talk")
		if err != nil{
			fmt.Println("talk1 faided",err)
			return
		}
		c.JSON(http.StatusOK, talk)
		fmt.Println("select:",talk)
	})

	r.GET("/tlakLike",func(c *gin.Context){
		db := mysqlInit()
		defer db.Close()
		var vnum int
		var lnum int
		err := db.QueryRow("select viewsNumber,likeNum from talk where id = ? " ,c.Query("tlakId")).Scan(&vnum,&lnum)
		db.Exec("update talk set viewsNumber = ?,likeNum = ?  where id = ?", vnum+1,lnum+1, c.Query("tlakId"))
		if err != nil{
			fmt.Println("exec faided",err)
			return
		}
		vnum = vnum + 1
		str := strconv.Itoa(vnum)
		lnum = lnum + 1
		str1 := strconv.Itoa(lnum)
		c.JSON(http.StatusOK, gin.H{
			"vnum" : str,
			"lnum" : str1,
			"message": "yes",
		})
	})
	r.GET("/viewNum",func(c *gin.Context){
		db := mysqlInit()
		defer db.Close()
		var vnum int
		err := db.QueryRow("select viewsNumber from talk where id = ? " ,c.Query("tlakId")).Scan(&vnum)
		db.Exec("update talk set viewsNumber = ? where id = ?", vnum+1, c.Query("tlakId"))
		if err != nil{
			fmt.Println("exec faided",err)
			return
		}
		vnum = vnum + 1
		str := strconv.Itoa(vnum)
		c.JSON(http.StatusOK, gin.H{
			"vnum" : str,
			"message": "yes",
		})
	})
	r.GET("/commentNum",func(c *gin.Context){
		db := mysqlInit()
		defer db.Close()
		var cnum int
		err := db.QueryRow("select comment from talk where id = ? " ,c.Query("tlakId")).Scan(&cnum)
		db.Exec("update talk set comment = ? where id = ?", cnum+1, c.Query("tlakId"))
		if err != nil{
			fmt.Println("exec faided",err)
			return
		}
		cnum = cnum + 1
		str := strconv.Itoa(cnum)
		c.JSON(http.StatusOK, gin.H{
			"cnum" : str,
			"message": "yes",
		})
	})

	r.GET("/talkComment",func(c *gin.Context){
		db := mysqlInitx()
		defer db.Close()
		var comments []TalkComment
		fmt.Println(c.Query("talkId"))
		err := db.Select(&comments,"select * from comment where talkId = ?",c.Query("talkId"))
		if err != nil{
			fmt.Println("talkComment faided",err)
			return
		}
		c.JSON(http.StatusOK, comments)
		fmt.Println("comments:",comments)
	})
	r.GET("/sendComment1",func(c *gin.Context){
		db := mysqlInit()
		defer db.Close()
		stmt, err := db.Prepare("insert into comment set replyId = ?,talkId = ?,positionFlag = ?,commentName = ?,objName = ?,QQ = ?,commentText = ?,commentTime = ?")
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := stmt.Exec(0,c.Query("talkId"),0,c.Query("commentName"),"Aurora",c.Query("qq"),c.Query("commentText"),c.Query("commentTime"))
		id, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println(id)
		c.JSON(http.StatusOK, gin.H{
			"message": "yes",
		})
	})
	r.GET("/sendComment2",func(c *gin.Context){
		db := mysqlInitx()
		defer db.Close()
		fmt.Println(c.Query("replyId"))
		r, err := db.Exec("insert into comment(talkId,replyId,positionFlag,commentName, objName, QQ,commentText, commentTime)values(?, ?, ?, ?, ?, ?, ?, ?)", c.Query("talkId"), c.Query("replyId"),  "1", c.Query("commentName"), c.Query("objName"), c.Query("qq"),c.Query("commentText") ,c.Query("commentTime"))
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		id, err := r.LastInsertId()
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		fmt.Println("insert succ:", id)
		c.JSON(http.StatusOK, gin.H{
			"message": "yes",
		})
	})
	r.GET("/sendComment3",func(c *gin.Context){
		db := mysqlInitx()
		defer db.Close()
		fmt.Println(c.Query("replyId"))
		r, err := db.Exec("insert into comment(talkId,replyId,positionFlag,commentName, objName, QQ,commentText, commentTime)values(?, ?, ?, ?, ?, ?, ?, ?)", c.Query("talkId"), c.Query("replyId"),  "2", c.Query("commentName"), c.Query("objName"), c.Query("qq"),c.Query("commentText") ,c.Query("commentTime"))
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		id, err := r.LastInsertId()
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		fmt.Println("insert succ:", id)
		c.JSON(http.StatusOK, gin.H{
			"message": "yes",
		})
	})
	r.GET("/algorithm",func(c *gin.Context){
		c.HTML(http.StatusOK, "算法.html", gin.H{
			"message": "yes",
		})
	})

	r.GET("/deepLearn",func(c *gin.Context){
		c.HTML(http.StatusOK, "深度学习.html", gin.H{
			"message": "yes",
		})
	})
	r.GET("/bulletScreen",func(c *gin.Context){
		c.HTML(http.StatusOK, "bulletScreen.html", gin.H{
			"message": "yes",
			"name": c.Query("name"),
			"qq": c.Query("qq"),
		})
	})
	r.GET("/bulletScreenMysql",func(c *gin.Context){
		db := mysqlInitx()
		defer db.Close()

		r, err := db.Exec("insert into bulletScreen(qq,userName,bulletScreenText)values(?, ?, ?)", c.Query("qq"), c.Query("userName"),  c.Query("bulletScreenText"))
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		id, err := r.LastInsertId()
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		fmt.Println("insert succ:", id)
		c.JSON(http.StatusOK, gin.H{
			"message": "yes",
		})
	})

	r.GET("/bulletScreenMysqldata",func(c *gin.Context){
		db := mysqlInitx()
		defer db.Close()
		var bulletScreenS []bulletScreen
		err := db.Select(&bulletScreenS,"select * from bulletScreen ")
		if err != nil{
			fmt.Println("bulletScreenMysqldata faided",err)
			return
		}
		c.JSON(http.StatusOK, bulletScreenS)
		fmt.Println("bulletScreenS:",bulletScreenS)
	})

	r.GET("/bulletScreenLatest",func(c *gin.Context){
		db := mysqlInitx()
		defer db.Close()
		var bulletScreenS []bulletScreen
		err := db.Select(&bulletScreenS,"select * from bulletScreen order by id desc limit 10")
		if err != nil{
			fmt.Println("bulletScreenMysqldata faided",err)
			return
		}
		c.JSON(http.StatusOK, bulletScreenS)
		fmt.Println("bulletScreenS:",bulletScreenS)
	})

	r.Run(":80")
}
func mysqlInit() *sql.DB{
	db, err := sql.Open("mysql", "root:XXXXXX@tcp(47.94.145.XXX:3306)/myweb?charset=utf8")
	if err != nil {
		panic(err)
	}
	return db
}
func mysqlInitx() *sqlx.DB{
	db, err := sqlx.Open("mysql", "root:XXXXX@tcp(47.94.145.XXX:3306)/myweb?charset=utf8")
	if err != nil {
		panic(err)
	}
	return db
}
type Email struct {
	From     string
	Host     string
	Port     int
	UserName string
	Password string
}
func (r *Email) Send(subject, htmlBody, attachFile, rename string, to ...string) error {
	l := len(to)
	if l == 0 {
		return errors.New("to can not be empty")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", r.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)
	if utils.FileIsExist(attachFile) {
		if rename != "" {
			setting := gomail.Rename(rename)
			m.Attach(attachFile, setting)
		} else {
			m.Attach(attachFile)
		}
	}
	d := gomail.NewDialer(r.Host, r.Port, r.UserName, r.Password)
	return d.DialAndSend(m)
}
func sendEmail(email string, str int) {
	var em = &Email{
		From:     "100354@qq.com",
		Host:     "smp.qq.com",
		Port:     45, //使用SSL，端口号465或587
		UserName: "11354",
		Password: "bxwzhyvbdgg",
	}
	hb := "<h3>" + strconv.Itoa(str) + "</h3>"                          //支持html
	err := em.Send("golang", hb, "", "", email) //可以多个接收方
	fmt.Println(err)
}