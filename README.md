GO-ZhihuDaily
=============

#[知乎日报 Web版（GoLang实现）](http://zhihudaily.ahorn.me/)

写这个的最初目的是因为PC上没有一个很好的阅读知乎日报的方式    
即那种一眼望过去，几天内的内容尽收眼底，然后点开自己感兴趣的Title继续阅读  

`cd GO-ZhihuDaily`  
`go run  main.go`  

浏览器url:   
`0.0.0.0:8000`   

---

API（点击看大图）:

![](https://github.com/Artwalk/GO-ZhihuDaily/blob/master/API.png)


---

**License**

GNU GENERAL PUBLIC LICENSE

---
以前没做过Web开发，边写边学GoLang/Git/HTML/CSS/GAE

( ⊙o⊙ )哇，这么多 '/'  
弱爆了有木有

域名、VPS都是蹭朋友的
太惨了

( >﹏<。)～呜呜呜……


---

1. 感谢 @[faceair](https://github.com/faceair/zhihudaily)，他做了最早的web版(PHP)（API就是从他的代码里找到的）
2. SQLite存储API返回的JSON数据，减小访问官网次数
3. 每小时更新一次当天数据
4. [Martini](https://github.com/codegangsta/martini) 框架
5. 蹭朋友的VPS [貌似十分不稳定，动不动就502了，（好吧，是我的小程序不稳定）]  

---

2014-02-24

图省事先用图片代替了

后面考虑用文字，这样复制粘贴也容易点

好吧，是我HTML/CSS不会啊，写起来步步维艰，现在真心做不到(>_<)

希望做成的样纸是 OldReader 那样，左边一栏标题，右边内容，实现滚动阅读

---


2014-02-25

V友太凶残了，先是502，下午论坛上发现后去SSH，发现进程已经没了，重开了一次   
晚上看htop，发现又已经吃掉VPS 25%的内存（1G），然后小伙伴生气啦(๑′°︿°๑)   

赶紧找BUG，这里改改那里动动，内存增速居然放缓了，赶紧先布上去再说   

然后慢慢改，好像是查询数据库后没close，还有上次更新时间忘重新赋值了，还有...，还有... （┬＿┬）   
之后剥离判断，加了个自动更新    

好不容易写个东东，那么简单还冒出来这么多问题，脸丢完了都╮(╯▽╰)╭

还好现在稳定了，据观察，一段时间后，内存不增反减  

GC好神奇

---

2014-02-28

貌似内存已经稳定了

请轻拍

---

2014-03-04  

我去，好像昨天图片全挂了，开始以为网络有问题  
后来发现好像被知乎日报官方屏蔽了，图片`403 forbidden`

只有先换成了文字顶上  

小伙伴伸出了援助之手，说图片缓存到他的VPS上好了  

（他是不知道啊，为啥我之前酱紫做呢？一共6W多条记录，所有图片下下来高达2G多）  

既然他说了，那就放上去好了  ( ≖‿≖)   

加了个下载，递归二分，8个狗同时跑

---

2014-03-05

用`convert crop`了图片头部，缓存不到200M，载入页面也快了

PS: 后台有一些奇怪的数据，比如没访问主页`/`，直接跳到`/page/62`，还有一些老的跳转`/url/**`，是因为浏览器缓存吗？统统重定向到主页`/`去了

---

2014-03-09

直接运行是只能下载裁剪当天的图片的，要在另一个地方调用  
理论上我应该加个注释，然后重构精简下，再加个图片下载成功判断   
但人都有惰性的，能跑就想着有时间再来好好优化  
最最重要的原因是，最近失恋了( ＞﹏＜)，对不住fork的大家了  


---

2014-03-26   

也没人issues我，知乎上发现有人评论才发现 又 502 了    
大惊，肿么又挂了啊...   

找BUG的过程中不停的想起《代码大全》里的一句话：  

>写的时候只有我和上帝知道是什么意思，而现在**只有上帝知道了**  

还好最后发现是日爆API改了，没有分享的图片了，招致臭名昭著的**空指针**  

话说这么老早的API了，还变JSON字段，还让不让我等P民活了丫   
~~~~(>_<)~~~~ 
