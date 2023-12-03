# Resources

Here are some prompts and resources I used to build this app:

## Links

- [HTML rendering in Go Gin framework](https://gin-gonic.com/docs/examples/html-rendering/)
- [JWT Authentication in Go](https://www.youtube.com/watch?v=ma7rUS_vW9M)
- [Example Go Gin CRUD](https://github.com/lemoncode21/golang-crud-gin-gorm)
- [Go Gin handling form input](https://gin-gonic.com/docs/examples/multipart-urlencoded-form/)

## Prompts

I have selected the most interesting / productive prompts, where with relatively little provided info I already received rather quality answers.
However, of course, several tweaks and improvements had to be made.

**1.**

Write me the beers_detail.tmpl file which would take this Beer as a param and display it.

```go
type BeerDetailed struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Thumbnail    string    `json:"thumbnail"`
	CommentCount uint      `json:"comment_count"`
	Comments     []Comment `json:"comments"`
	UpvoteCount  uint      `json:"upvote_count"`
	AuthorId     uint      `json:"author_id"`
}

type Comment struct {
	ID          uint   `json:"id"`
	AuthorID    uint   `json:"user"`
	Author      string `json:"user_name"`
	Content     string `json:"content"`
	CreatedDate int64  `json:"created_date"`
	BeerID      uint   `json:"beer"`
}
```

<hr>

**2.**

I have a following auth login handler in my go gin web app:
```go
func (ctrl *UserController) LoginHandler(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse request body",
		})
		return
	}
	user := ctrl.authService.Login(body.Email, body.Password)

	// generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// save to cookies
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenStr, 3600*24*30, "", "", false, true)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{})
}

```
Write me a handler for logout and show how to add it to router

<hr>

**3.**

Write me an index page in plain html and css that would be like a landing page for a beer website. This landing page must have a header with button to go to beer list. In the main page part, there should be a couple of sections with images and some related stuff + one more button that would bring user to the beer list

<hr>

**4.**

I have a following line of code to retrieve a row from database in golang:
```go
row := db.QueryRow("SELECT * FROM beers WHERE ID = $1", beerID)
```
How do I check if it has returned some value or not?