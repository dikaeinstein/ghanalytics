package analytics_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/dikaeinstein/ghanalytics/analytics"
	"github.com/dikaeinstein/ghanalytics/data"
)

func TestListTopNUsers(t *testing.T) {
	expectedActors := []analytics.Actor{
		{8517910, "LombiqBot"},
		{29139614, "renovate[bot]"},
		{44826218, "MatoPlus"},
		{625469, "armano2"},
		{1008205, "wigforss"},
		{2134633, "buddyspike"},
		{2631623, "onosendi"},
		{2895902, "patsonluk"},
		{5271692, "romankagan"},
		{5954907, "awesomekling"},
		{2631623, "onosendi"},
		{24802628, "gpeterson406"},
		{37038626, "geos4s"},
		{46163555, "PySimpleGUI"},
		{2895902, "patsonluk"},
		{30060991, "m41na"},
		{43908490, "jick155"},
	}

	testCases := []struct {
		desc          string
		limit         int
		sortCriterion []analytics.SortCriteria
	}{
		{
			desc:  "Top 5 Users",
			limit: 5,
			sortCriterion: []analytics.SortCriteria{
				analytics.CommitsPushed,
				analytics.PrCreated,
			},
		},
		{
			desc:  "Top 10 Users",
			limit: 10,
			sortCriterion: []analytics.SortCriteria{
				analytics.CommitsPushed,
				analytics.PrCreated,
			},
		},
	}

	store := createStore(t)
	a := analytics.New(store)

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			topNUsers, err := a.ListUsers(
				analytics.Sort(tC.sortCriterion),
				analytics.Limit(tC.limit),
			)
			if err != nil {
				t.Error(err)
			}

			if len(topNUsers) != tC.limit {
				t.Errorf("Length of topNUsers does not match. want %d; got %d",
					tC.limit, len(topNUsers))
			}

			expectedNActors := expectedActors[0:tC.limit]
			if !reflect.DeepEqual(topNUsers, expectedNActors) {
				t.Errorf("Wrong topNUsers returned. want %+v; got %+v",
					expectedNActors, topNUsers)
			}
		})
	}
}

func TestListTopNRepos(t *testing.T) {
	expectedActors := []analytics.Repo{
		{ID: 42018768, Name: "Lombiq/Helpful-Libraries"},
		{ID: 62069489, Name: "multicharts/scanner-check"},
		{ID: 204268723, Name: "supershell2019/conf"},
		{ID: 231160326, Name: "MatoPlus/react-jsx-lab-cb-gh-000"},
		{ID: 6495132, Name: "wikimedia/mediawiki-extensions-CentralAuth"},
		{ID: 20598648, Name: "ElvUI-WotLK/ElvUI"},
		{ID: 45993266, Name: "patsonluk/airline"},
		{ID: 53678002, Name: "qiyanjun/books2read"},
		{ID: 67651864, Name: "Lombiq/Git-Hg-Mirror-Common"},
		{ID: 83585690, Name: "TedTschopp/tschopp.net"},
	}

	testCases := []struct {
		desc          string
		limit         int
		sortCriterion []analytics.SortCriteria
	}{
		{
			desc:  "Top 5 Repos",
			limit: 5,
			sortCriterion: []analytics.SortCriteria{
				analytics.CommitsPushed,
			},
		},
		{
			desc:  "Top 10 Repos",
			limit: 10,
			sortCriterion: []analytics.SortCriteria{
				analytics.CommitsPushed,
			},
		},
	}

	store := createStore(t)
	a := analytics.New(store)

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			topNRepos, err := a.ListRepos(
				analytics.Sort(tC.sortCriterion),
				analytics.Limit(tC.limit),
			)
			if err != nil {
				t.Error(err)
			}

			if len(topNRepos) != tC.limit {
				t.Errorf("Length of topNRepos does not match. want %d; got %d",
					tC.limit, len(topNRepos))
			}

			expectedNActors := expectedActors[0:tC.limit]
			if !reflect.DeepEqual(topNRepos, expectedNActors) {
				t.Errorf("Wrong topNRepos returned. want %+v; got %+v",
					expectedNActors, topNRepos)
			}
		})
	}
}

func createStore(t *testing.T) *data.Store {
	t.Helper()

	store, err := data.NewStore(actorsCSVFile(t), commitsCSVFile(t),
		eventsCSVFile(t), reposCSVFile(t))
	if err != nil {
		t.Fatal(err)
	}

	return store
}

func actorsCSVFile(t *testing.T) io.Reader {
	t.Helper()

	return strings.NewReader(
		`id,username
8422699,Apexal
53201765,ArturoCamacho0
2631623,onosendi
52553915,anggi1234
31390726,AdrianWilczynski
5954907,awesomekling
10052381,PercussiveElbow
30060991,m41na
8517910,LombiqBot
19255925,nachoapps
2254036,mpro34
2134633,buddyspike
48773997,Danver98
24802628,gpeterson406
44826218,MatoPlus
13025337,bitnami-bot
1008205,wigforss
8517910,LombiqBot
34975776,ryo97321
5271692,romankagan
44267780,thadiun
52545156,maromaroXD
625469,armano2
46163555,PySimpleGUI
37038626,geos4s
59149927,AssamToumi
1686007,efarbereger
2895902,patsonluk
43908490,jick155
29139614,renovate[bot]
6442743,dvldrmmr
29139614,renovate[bot]
13025337,bitnami-bot
26935674,rodrigoferraro
2733465,yubo
44826218,MatoPlus
8517910,LombiqBot`,
	)
}

func commitsCSVFile(t *testing.T) io.Reader {
	t.Helper()

	return strings.NewReader(`sha,message,event_id
5948a6cc5255015e983a9719117c15ff197b4681,Refactor member inde,11185376329
bf7296401598660b44d8923787a2600f346f9a81,Refactor roadmap,11185376329
488794042fce073c5075180becc9bfaf1156eb7e,Member detail start,11185376329
18fe8470f986fb009ee1f6841546382416167d3f,Add review to Event model,11185376329
2a7dda4b06bb9762664eaf5493ce674c0cb2a113,Remove masaka - not being used,11185376335
9784df251c4aafbb322cac284f92b4b377d5b78e,Add files via upload,11185376336
28fcf3a7c1ed3bd182e7aa030bc37b92ad58c2d9,Save README in specified location (CLI argument). Append new content to previous content.,11185376338
8602fa5b49aa4e2b039764a14698f0baa3ad0532,"Kernel: Enable x86 SMEP (Supervisor Mode Execution Protection)

This prevents the kernel from jumping to code in userspace memory.",11185376339
1e94e09574d9a88a704ebefd3b0bb2edd1a56473,added resilience decorators to the http client and cleaned up a bunch of useless artifacts,11185376342
597d371b38506813e0766ae3e2453493884efaad,"Merge pull request #3 from m41na/v3-apply-resilience

added resilience decorators to the http client and cleaned up a bunch…",11185376342
9f74d7ce6d136ac5f3749f017f694fdb9d2b018d,data updates,11185376344
2636e643f4479b686e6309f6845d8a6101e65897,Commit command,11185376352
30f6f279cc03fe4976a52fbd5269c4d09a09bf44,checking for request object 2,11185376356
0766285421fe035bffbebb8f92bb1cae4bd578cd,Avoid using xcolor due to conflicts,11185376357
22462c0b322d8b643595db975090b29e0c0536fd,Use figures with dash instead of dot names,11185376357
f55231a695d42910d93618aa20654b68ba929a93,Chapter 2 and 3 layout fixes,11185376357
1d6c4e79774a20fefcf533b17d9c56ecff9f6865,PDF with updates through chapter 3.,11185376357
56fe953733f4acb9349659e775ccc4e796e20ec7,Automatically backed up by Learn,11185376360
0aaf1a059304cbf843c670338f75ba2685d86a99,1.10.0-debian-9-r11 release,11185376363
707b45099067ff981098ebd9a50c3df15fa3a91e,add eyepatch/,11185376371
1da920a1e9197f54a03005a887ff44098f4ea249,"Changed theme's ""black"" and ""white"" values to hex values.",11185376380
bb70955d4677ec753be7ecfc9f21cd37c43800b4,Add sgroup07/metar_edhi.csv,11185376381
960300dacb59bb10bb008afd6663f5d36a9b3d54,Optimize airplane home patching,11185376388
95aa094c8e0211b3471949e0213b905366aa4bd8,add item to readed me,11185376392
2d682f937d44cebc1df63667750057cee50810fe,Remove embedded bootstrap and jquery in favor of CDN,11185376393
38d1ed56ee9943ede97a4c4eee39becb15b2648e,add user image,11185376393
d11022dba46075915cd87e8fa7ba29519ea651b3,more admin,11185376393
effbb106c69637e796fa95d397a0ac99968e0ec3,remove unnecessary css files,11185376393
7ccd36df8f080915463b8d2a2f66035f43f040b1,Add renovate.json,11185376393
32891ad1b2ed1ff25bbad0ec8f784ec52a3a4360,uploading social interaction,11185376402
badd55aaab8f548f11d80b86b6035118ab22a1bb,Add renovate.json,11185376402
133ed7ec0f6c93b68048e8faa8eb50b14ad9051a,List provider notification,11185376409
bc4e9c9d1b6bb24ab06e1401e56da198a9267dc8,Automatically backed up by Learn,11185376411`,
	)
}

func eventsCSVFile(t *testing.T) io.Reader {
	t.Helper()

	return strings.NewReader(`id,type,actor_id,repo_id
11185376329,PushEvent,8422699,224252202
11185376333,CreateEvent,53201765,231161852
11185376335,PushEvent,2631623,155254893
11185376336,PushEvent,52553915,231065965
11185376338,PushEvent,31390726,225080339
11185376339,PushEvent,5954907,160083795
11185376341,WatchEvent,10052381,221552739
11185376342,PushEvent,30060991,230923653
11185376343,PushEvent,8517910,107471694
11185376344,PushEvent,19255925,223831715
11185376350,CreateEvent,2254036,231161853
11185376352,PushEvent,2134633,231031280
11185376356,PushEvent,48773997,227178623
11185376357,PushEvent,24802628,100496729
11185376360,PushEvent,44826218,231160326
11185376363,PushEvent,13025337,194280481
11185376364,CreateEvent,1008205,230448690
11185376369,PushEvent,8517910,42018768
11185376371,PushEvent,34975776,197711881
11185376374,PullRequestEvent,5271692,138499772
11185376375,IssuesEvent,44267780,153640774
11185376376,DeleteEvent,52545156,211517583
11185376377,PullRequestEvent,625469,165536154
11185376380,PushEvent,46163555,140614233
11185376381,PushEvent,37038626,159298221
11185376383,CreateEvent,59149927,231161840
11185376387,CreateEvent,1686007,231161855
11185376388,PushEvent,2895902,45993266
11185376392,PushEvent,43908490,228343155
11185376393,PushEvent,29139614,230275330
11185376399,WatchEvent,6442743,153830667
11185376402,PushEvent,29139614,83585690
11185376406,CreateEvent,13025337,194280481
11185376409,PushEvent,26935674,228284912
11185376410,DeleteEvent,2733465,11581991
11185376411,PushEvent,44826218,231160326
11185376421,PushEvent,8517910,42018768
11185376423,PushEvent,5699167,53678002
11185376425,PushEvent,32727560,230797706
11185376428,CreateEvent,30938,231161498
11185376430,IssuesEvent,27738369,207424413
11185376432,PushEvent,2539292,6495132
11185376434,PushEvent,2797449,108752559
11185376436,PushEvent,54496419,204268723
11185376439,PushEvent,3224548,231161702
11185376442,PushEvent,38356640,129512184
11185376443,WatchEvent,9023479,31792824
11185376445,PushEvent,32480,231161718
11185376446,PushEvent,5816685,20598648
11185376447,PushEvent,1008205,230448690
11185376452,PushEvent,8676741,170661918
11185376453,IssueCommentEvent,53356952,165536154
11185376455,PushEvent,17928155,62069489
11185376456,CreateEvent,17261190,190052895
11185376457,PushEvent,40586421,138676186
11185376458,PushEvent,48690212,212703717
11185376469,CreateEvent,51274871,231161859
11185376473,IssueCommentEvent,12852815,108332615
11185376474,PushEvent,42107416,231160806
11185376480,PushEvent,8517910,67651864
11185376487,PushEvent,4282702,173629625
11185376488,PushEvent,45283053,198916487
11185376485,PushEvent,45327612,159010660
11185376491,PushEvent,30926704,231158291
11185376492,PushEvent,5253690,153466799
11185376494,PushEvent,17928155,62069489
11185376498,PushEvent,2486411,135042737
11185376505,PushEvent,541490,224001345
11185376508,PushEvent,54496419,204268723
11185376511,ForkEvent,3623783,230809890
11185376512,PushEvent,6701060,84128241
11185376514,IssueCommentEvent,54035,44727837
11185376516,PushEvent,48352006,229235126
11185376519,PushEvent,29507018,136921024
11185376528,PushEvent,40587912,138681984
11185376531,WatchEvent,16427717,218329977`,
	)
}

func reposCSVFile(t *testing.T) io.Reader {
	t.Helper()

	return strings.NewReader(`id,name
224252202,DSC-RPI/dsc-portal
231161852,ArturoCamacho0/ProjectResponsive
155254893,onosendi/flask-dlindegren
231065965,anggi1234/tunas-ilmu
225080339,AdrianWilczynski/VSCodeSnippetsToVS
160083795,SerenityOS/serenity
221552739,z0ph/aws-security-toolbox
230923653,m41na/demo-micro-services
107471694,Lombiq/PJS.ReTouch
223831715,nachoapps/pad-game-data
231161853,mpro34/GunGraves
231031280,buddyspike/gogit
227178623,Danver98/sneakers-backend
100496729,gpeterson406/Greenwood_Book
231160326,MatoPlus/react-jsx-lab-cb-gh-000
194280481,bitnami/bitnami-docker-harbor-portal
230448690,wigforss/java-8-base
42018768,Lombiq/Helpful-Libraries
197711881,ryo97321/rust-workspace
138499772,nazmulidris/algorithms-in-kotlin
153640774,thadiun/hello-world
211517583,maromaroXD/FiFO_Second_Chance
165536154,typescript-eslint/typescript-eslint
140614233,PySimpleGUI/PySimpleGUI
159298221,geos4s/18w856162
231161840,AssamToumi/Bootstrap-blog
231161855,efarbereger/tmp_clock_repo
45993266,patsonluk/airline
228343155,jick155/Machine-Learning-portfolio
230275330,mrlynn/mongodb-github-student
153830667,KrzysztofSzewczyk/asmbf
83585690,TedTschopp/tschopp.net
194280481,bitnami/bitnami-docker-harbor-portal
228284912,rodrigoferraro/bootcamp2019
11581991,yubo/doc
231160326,MatoPlus/react-jsx-lab-cb-gh-000
42018768,Lombiq/Helpful-Libraries
53678002,qiyanjun/books2read
230797706,BlackHoleSecurity/website
231161498,sullis/amazon-ecr-login
207424413,skrapkam/samchang
6495132,wikimedia/mediawiki-extensions-CentralAuth
108752559,BenedekFarkas/Elemental-Gankery
204268723,supershell2019/conf
231161702,fatalhalt/LEDBlink
129512184,SockPuppetry/Loop-A
31792824,flutter/flutter
231161718,jethrolarson/parse-ts
20598648,ElvUI-WotLK/ElvUI
230448690,wigforss/java-8-base
170661918,nulastudio/Freedom
165536154,typescript-eslint/typescript-eslint
62069489,multicharts/scanner-check
190052895,ran-dall/homebrew-cask-drivers
138676186,himobi/hotspot
212703717,DemiMarie-parity/rust-libp2p
231161859,konstantinkot/portfolio
108332615,SalesforceFoundation/gem
231160806,kilcmn3/programming-univbasics-nds-hashes-of-hashes-lab-nyc-web-010620
67651864,Lombiq/Git-Hg-Mirror-Common
173629625,thijs852/parkeer_data
198916487,flycnzebra/itop-4412-android5.1
159010660,0x21fr2/Hacking
231158291,haitdhn/Petstore`,
	)
}
