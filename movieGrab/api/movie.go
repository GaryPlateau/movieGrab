package api

type MovieWeb struct {
	WebName string
	Types   []MovieType
}

type MovieType struct {
	TypeId    int
	TypeName  string
	MovieKind string
}

type UrlType struct {
	UrlId   int
	UrlName string
	UrlLink string
}

func GetDygangMovieType() []*MovieType {
	return []*MovieType{
		{1, "恐怖片", "kongbupian"},
		{2, "喜剧片", "xijupian"},
		{3, "动作片", "dongzuopian"},
		{4, "爱情片", "aiqingpian"},
		{5, "科幻片", "kehuanpian"},
		{6, "战争片", "zhanzhengpian"},
		{7, "悬疑片", "xuanyipian"},
	}
}

func GetDytt8MovieType() []*MovieType {
	return []*MovieType{
		{1, "动作片", "dongzuopian"},
		{2, "剧情片", "juqingpian"},
		{3, "爱情片", "aiqingpian"},
		{4, "喜剧片", "xijupian"},
		{5, "科幻片", "kehuanpian"},
		{6, "恐怖片", "kongbupian"},
		{7, "动画片", "donghuapian"},
		{8, "惊悚片", "jingsongpian"},
		{9, "战争片", "zhanzhengpian"},
		{10, "犯罪片", "fanzuipian"},
		{11, "灾难片", "zainanpian"},
		{12, "纪录片", "jilupian"},
		{13, "奇幻片", "qihuanpian"},
	}
}

func Get66ysMovieType() []*MovieType {
	return []*MovieType{
		{1, "动作片", "dongzuopian"},
		{2, "恐怖片", "kongbupian"},
		{3, "战争片", "zhanzhengpian"},
		{4, "科幻片", "kehuanpian"},
		{5, "爱情片", "aiqingpian"},
		{6, "喜剧片", "xijupian"},
		{7, "纪录片", "jilupian"},
		{8, "剧情片", "bd"},
		{9, "国产剧", "dsj"},
		{10, "港台剧", "dsj2"},
		{11, "日韩剧", "dsj3"},
		{12, "欧美剧", "dsj4"},
		{13, "国配电影", "gy"},
	}
}

func GetAllMovieType() []*MovieType {
	return []*MovieType{
		{0, "全网搜索", "allMovie"},
	}
}
