package main

import (
	"encoding/json"
	"net/http"
)

type PostItem struct {
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Img        string `json:"img"`
	AuthorLogo string `json:"authorLogo"`
	AuthorName string `json:"authorName"`
	Date       string `json:"date"`
	Href       string `json:"href"`
}
type Posts []PostItem

var posts = []PostItem{
	{
		Title:      "祝你 1024 节日快乐",
		Desc:       "🎉今天是 1024 程序员节，平淡的日子需要有些仪式感，周末写了一个小彩蛋，送给程序员以及喜欢开源的同学，希望你喜欢，也祝福你以后写的都是你想写的代码。由于没有任何署名，当然你分享给你的程序员朋友也没有问题的。",
		Img:        "https://tw93.netlify.app/img/1024-banner.png",
		AuthorLogo: "https://pbs.twimg.com/profile_images/1540397753586528256/SFkyn7LD_400x400.jpg",
		AuthorName: "Tw93",
		Date:       "2022-10-24",
		Href:       "https://tw93.netlify.app/",
	},
	{
		Title:      "Notion 静态博客生成器",
		Desc:       "一个使用 NextJS + Notion API 实现的，部署在 Vercel 上的静态博客系统，只要在 Notion 写好文章就会自动同发布为静态博客，从而专注于写作、而不需要操心网站的维护。",
		Img:        "https://pbs.twimg.com/media/Ffvv4zRVQAECrCS?format=jpg&name=large",
		AuthorLogo: "https://pbs.twimg.com/profile_images/1477903349278453762/0OBeufkj_400x400.jpg",
		AuthorName: "龙爪槐守望者",
		Date:       "2022-10-23",
		Href:       "https://github.com/tangly1024/NotionNext",
	},
	{
		Title:      "突然发现，其实与其说喜欢折腾，不如说多巴胺成瘾 - V2EX",
		Desc:       "引起很多人共鸣，俺也一样。🥹",
		Img:        "https://pbs.twimg.com/media/FfP6nq_UcAASu6v?format=png&name=medium",
		AuthorLogo: "https://pbs.twimg.com/profile_images/1548486285970444289/rALdZ0SF_400x400.png",
		AuthorName: "野生架构师 🐒",
		Date:       "2022-10-17",
		Href:       "https://www.v2ex.com/t/887301",
	},
	{
		Title:      "How to Use Root C++ Interpreter Shell to Write C++ Programs",
		Desc:       "The powerful gravity waves resulting from the impact of the planets' moons â€” four in total â€” were finally resolved in 2015 when gravitational microlensing was used to observe the",
		Img:        "https://images.unsplash.com/photo-1617529497471-9218633199c0?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=870&q=80",
		AuthorLogo: "https://api.uifaces.co/our-content/donated/KtCFjlD4.jpg",
		AuthorName: "Lourin",
		Date:       "Jan 4 2022",
		Href:       "javascript:void(0)",
	},
	{
		Title:      "How to Use Root C++ Interpreter Shell to Write C++ Programs",
		Desc:       "The powerful gravity waves resulting from the impact of the planets' moons â€” four in total â€” were finally resolved in 2015 when gravitational microlensing was used to observe the",
		Img:        "https://images.unsplash.com/photo-1617529497471-9218633199c0?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=870&q=80",
		AuthorLogo: "https://api.uifaces.co/our-content/donated/KtCFjlD4.jpg",
		AuthorName: "Lourin",
		Date:       "Jan 4 2022",
		Href:       "javascript:void(0)",
	}, {
		Title:      "How to Use Root C++ Interpreter Shell to Write C++ Programs",
		Desc:       "The powerful gravity waves resulting from the impact of the planets' moons â€” four in total â€” were finally resolved in 2015 when gravitational microlensing was used to observe the",
		Img:        "https://images.unsplash.com/photo-1617529497471-9218633199c0?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=870&q=80",
		AuthorLogo: "https://api.uifaces.co/our-content/donated/KtCFjlD4.jpg",
		AuthorName: "Lourin",
		Date:       "Jan 4 2022",
		Href:       "javascript:void(0)",
	}, {
		Title:      "How to Use Root C++ Interpreter Shell to Write C++ Programs",
		Desc:       "The powerful gravity waves resulting from the impact of the planets' moons â€” four in total â€” were finally resolved in 2015 when gravitational microlensing was used to observe the",
		Img:        "https://images.unsplash.com/photo-1617529497471-9218633199c0?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=870&q=80",
		AuthorLogo: "https://api.uifaces.co/our-content/donated/KtCFjlD4.jpg",
		AuthorName: "Lourin",
		Date:       "Jan 4 2022",
		Href:       "javascript:void(0)",
	}, {
		Title:      "How to Use Root C++ Interpreter Shell to Write C++ Programs",
		Desc:       "The powerful gravity waves resulting from the impact of the planets' moons â€” four in total â€” were finally resolved in 2015 when gravitational microlensing was used to observe the",
		Img:        "https://images.unsplash.com/photo-1617529497471-9218633199c0?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=870&q=80",
		AuthorLogo: "https://api.uifaces.co/our-content/donated/KtCFjlD4.jpg",
		AuthorName: "Lourin",
		Date:       "Jan 4 2022",
		Href:       "javascript:void(0)",
	}, {
		Title:      "How to Use Root C++ Interpreter Shell to Write C++ Programs",
		Desc:       "The powerful gravity waves resulting from the impact of the planets' moons â€” four in total â€” were finally resolved in 2015 when gravitational microlensing was used to observe the",
		Img:        "https://images.unsplash.com/photo-1617529497471-9218633199c0?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=870&q=80",
		AuthorLogo: "https://api.uifaces.co/our-content/donated/KtCFjlD4.jpg",
		AuthorName: "Lourin",
		Date:       "Jan 4 2022",
		Href:       "javascript:void(0)",
	}, {
		Title:      "How to Use Root C++ Interpreter Shell to Write C++ Programs",
		Desc:       "The powerful gravity waves resulting from the impact of the planets' moons â€” four in total â€” were finally resolved in 2015 when gravitational microlensing was used to observe the",
		Img:        "https://images.unsplash.com/photo-1617529497471-9218633199c0?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=870&q=80",
		AuthorLogo: "https://api.uifaces.co/our-content/donated/KtCFjlD4.jpg",
		AuthorName: "Lourin",
		Date:       "Jan 4 2022",
		Href:       "javascript:void(0)",
	}, {
		Title:      "How to Use Root C++ Interpreter Shell to Write C++ Programs",
		Desc:       "The powerful gravity waves resulting from the impact of the planets' moons â€” four in total â€” were finally resolved in 2015 when gravitational microlensing was used to observe the",
		Img:        "https://images.unsplash.com/photo-1617529497471-9218633199c0?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=870&q=80",
		AuthorLogo: "https://api.uifaces.co/our-content/donated/KtCFjlD4.jpg",
		AuthorName: "Lourin",
		Date:       "Jan 4 2022",
		Href:       "javascript:void(0)",
	},
}

func main() {

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {

		postsJSON, _ := json.Marshal(posts)
		w.Header().Set("Content-Type", "application/json")
		w.Write(postsJSON)

	})
	http.ListenAndServe(":4000", nil)
}
