# Todo-App
It is a simple todo list written in **golang** using local __**mysql**__ database, and __**html**__,__**css**__ and __**bootstrap**__ for front end.

## Application Requirement:
1. golang: https://golang.org/dl/<br>
2. mysql driver package to connect with MySQL: ```go get -u github.com/go-sql-driver/mysql```<br>

## Start the application
1.Create a .env file and copy the key for your db source to __.env.DB_SOURCE__ variable <br>
2.Start mysql server. <br>
3.Run ```go run main.go```  <br>
4.Open [server](http://localhost:8080/todo)

## Sample
![todo sample](todo.jpg)