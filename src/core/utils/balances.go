package balance

import (
	"fmt"
	"math/big"
	"portfolio/core/schema"
	"sort"

	Multicall "portfolio/core/contracts/MulticallContract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

type BalanceCall struct {
	tokenAddress  common.Address
	walletAddress common.Address
}

const BALANCE_OF_FUNC = "70a08231"

func BalanceOf(call BalanceCall) Multicall.Multicall3Call3 {
	// fmt.Println(call.tokenAddress, " - ", call.walletAddress.Hash().String()[2:], "   ")
	return Multicall.Multicall3Call3{
		Target:       call.tokenAddress,
		AllowFailure: true,
		CallData:     common.Hex2Bytes(fmt.Sprintf("%s%s", BALANCE_OF_FUNC, call.walletAddress.Hash().String()[2:]))}
}

func genGetMultipleBalanceCalls(tokens []common.Address, wallets []common.Address) []Multicall.Multicall3Call3 {
	res := make([]Multicall.Multicall3Call3, len(wallets)*len(tokens))
	counter := 0
	for _, token := range tokens {
		for _, wallet := range wallets {
			res[counter] = BalanceOf(BalanceCall{tokenAddress: token, walletAddress: wallet})
			counter++
		}
	}
	return res
}

type IndexedCall struct {
	index uint64
	call  Multicall.Multicall3Call3
}

func genGetBalanceCalls(tokens []schema.Token, wallet common.Address) []IndexedCall {
	res := make([]IndexedCall, len(tokens))
	var counter uint64 = 0
	for _, token := range tokens {
		res[counter] = IndexedCall{counter, BalanceOf(BalanceCall{tokenAddress: token.Address, walletAddress: wallet})}
		counter++
	}
	return res
}

func GetBalances(
	multiCaller Multicall.MulticallCaller,
	tokens []common.Address,
	wallets []common.Address) [][]big.Int {

	calls := genGetMultipleBalanceCalls(tokens, wallets)
	callRes, err := multiCaller.Aggregate3(&bind.CallOpts{}, calls)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	res := make([][]big.Int, len(wallets))
	var _res []big.Int
	i, j, tokenCount := 0, 0, len(tokens)

	for _, mr := range callRes {
		z := new(big.Int)
		z.SetBytes(mr.ReturnData)
		if j == 0 {
			_res = make([]big.Int, len(tokens))
		} else if j == tokenCount {
			j = 0
			res[i] = _res
		}
		i++
		j++
	}
	return res
}

func chunks[T any](array []T, chunkSize int) [][]T {
	steps := len(array) / chunkSize

	c := make([][]T, chunkSize)
	for i := 0; i < chunkSize; i++ {
		c[i] = array[(i * steps):((i + 1) * steps)]
	}
	return c
}

type chunkResult struct {
	index   uint64
	callRes Multicall.Multicall3Result
	err     any
}

func GetBalancesFaster(
	multiCaller Multicall.MulticallCaller,
	tokens []schema.Token,
	wallets common.Address) []schema.Token {

	allCalls := genGetBalanceCalls(tokens, wallets)

	chunkChannel := make(chan []chunkResult)

	chunkedCalls := chunks[IndexedCall](allCalls, 100)

	for _, indexCalls := range chunkedCalls {
		go func(indexCalls []IndexedCall, chunkChannel chan []chunkResult) {
			calls := make([]Multicall.Multicall3Call3, len(indexCalls))
			for i, indexedCall := range indexCalls {
				calls[i] = indexedCall.call
			}
			res, err := multiCaller.Aggregate3(&bind.CallOpts{}, calls)
			parsedRes := make([]chunkResult, len(indexCalls))
			for i, indexedCall := range indexCalls {
				parsedRes[i] = chunkResult{indexedCall.index, res[i], err}
			}
			chunkChannel <- parsedRes
		}(indexCalls, chunkChannel)
	}

	//for i := 0; i < len(tokens); i++ {
	for i := 0; i < len(chunkedCalls); i++ {

		_tmp := <-chunkChannel
		for _, tmp := range _tmp {

			if tmp.callRes.Success {
				z := new(big.Int)
				z.SetBytes(tmp.callRes.ReturnData)
				if len(z.Bits()) == 0 {
					continue
				}
				badNum := new(big.Int)
				// NOTE - some contracts returns this output instead of 0 :/
				// https://etherscan.io/address/0xe9Cf7887b93150D4F2Da7dFc6D502B216438F244#readProxyContract
				// 0xf977814e90da44bfa03b6295a0616a897441acec @ balanceOf

				badNum.SetString("5551322697364827663932540521968496641513208056639073144345637956117013275158518319394074513605658646672034377155866499211490207096881866628815899538692413447209068575243838503190958262344102506860710145468851577315649165604149695894808279194498659390011867486021885017898876353591592154411108381551525930841566196403539322915484446582580249514168108941008642981736048256454931837408132822660030072688995196886040162568544461224467619083068323177016887839072963784397575247245025359716506217923660935077004619478864426212158305452292955807777440157984054572364790366051734905727805768716718647125380705250974331925949739417426547496426910345365897533472356511240670865372206037352555714240312560926314559686920828887407736120106414203743255394129398102843834075700325063738618064740493105683475320751660534350135668912705031677400481341130357037170839675280675040183084436151635550460348329729277655151415450824131271456202226800099351690974137897792103341926423803370376505042388854638058027221947872634819215479016056460796559644880943524438265553742974055155295354377291653168014241825848862047992117396791707074632947579310697458439036519962231342020550266061217313721287955204430739262640545728451268867642844552956898287487338197980601674953321982869063185649872308076602367714783143741444663235237729476173148070683349112125064529283860349720556837837146362700570199476044270416325206315206578548140226116909954197745485763770437342498786974694734181794977563887398968202275475015858380485969706027843125012584014304754073114625510819364854446751564549050264741341735781882303255605996842930703534060937757480307432734884368611445398708870627925690524645692460523732607582725496363691626267018517326438149496028810882664034001622858179934864595345241372562457271804646634144676403074621792840784319182627343239748710670625037349908612501720758872688425676822149021791509641731412249405640945844983071041541102923132122000842468199401541961797593774150327139701213789732020221957179630541511499004457937095016068371688183475511272325269475790154852820333961050263445479688460344451574250421915378722354609846955892356310260688350902150017791267648195303826035956228762142043005915052261274866336322080843792191772047793115827072114886493219141485876023616500922657921056509596009708858619470805659919105329197871779163071359537936934216678881043549454836421653140534292850568896499899382380633489951766319925120743280730815695094692703971944307525402699824773795764231376682634303825471576778116646534695433556992634276065915083798277949474532766270149847002060436885556870638469241473357731266225809512256408506323502976460761135517794849344371796388907785058712218623501883681400042168953655433212498108868824680242290090793298779985347873818818716101178680513415688239918870443654418003859104497051918288254521049409590923631268671036919164159212361444683818241822249775979851392425002127684013823777318674668594346899981335955435705674383847397194234298808670829078221108861223038774425288508547927806825256653613817959172725241798338982989776973460975717729602357530799634340116301789817199175411605712285408181894750891121086164824875439287772263115519091807716582417418700299134451874768046876658908583061152720943422095889674755245092056829906985115190994008140744615573784682427701161574562936828714695837446399754520565855367125733253760348995172554765936279006863715661222911169821320137332931559615806918281180121197599613620690764272001315221347796893805433531909812276805923118612098479276590258007768487520580443557572604078201325811738267332732049888827009129520386562476175453655026223981384037655321130798712381594057457611199382791201055873488110852192216087357953764233262061367437198807181020188857354004149634793942544076324114343424436176522216199422958657501471087637900782805201865812360666272548702101209588157384124355836617993236749771940101362184949311694845315916690108727714840104015688394489341430944994955153340571112840139907844860619716656730489466355665499324359538747105234011359401489302688084969570396363779283876754608914395324435223570574400508422577377128672260056530052174881813204339721293391233784891439784011605761421212878288168856124273942615481139687169070077745755579517236492852665264108542852982400166696589647808894101017443625462275592986334309300292295654239257390129211744613392427528512382918417482127234331585770830831251559257046475948598105422264349425624040709552250368979115415442384641981999610067702028623467256495527302104470930808624070078410774976440841803155974429827004452891727752164455758552513616418458380893706037891293264398346147469888554259604093929770105635767590315326435172775010387900822772212597536577031013170890804105843340356155329228394239021671988127178245565480225816632373628338348528833910228156686843753840624910446911834952240810599743970040400204487610554478110064855026230261155260100757245542973901180309056320581873163286435013218153420986083206382522126724159504028436451117941729114686769211169995368659557340913602552308725018558983456103592127761167289188428649093680389474351076829502195863244765984150900755480105789723337317939817412282202459116897440238058572184532754033922230748808640940343613714524532564543227517495186912543952624805001616352652575573034725228210765624890395517121076434569919120475607038257472182308928271922933616225992035454379838169331348084831079129636736330615403915752464720965760404727751221385824396896375975269088125137436898549962868110913904352601999605050421452017597219149338642325662229236124556141721083506198294677252657943560158805156872569227499742560552355257497643787349849267225820392109741962888839161398568992588258824804377640769758665829839203292108531530906340337339521876412230461961195742109079162827373280431759016730254302790003396922266028016849746194582805906181308555659423226815979988416574981524160548673384129286054654230006943334015771545120552270836915358808810225112258001508247647743553081753381671500759657750077595872032084466998177155983460341548921365392484071774469497492932159065900102074162049632478064285549007031419667277281054703905789177186997862891694409414974767444316779710573385098941040104073916364556250855621070779596991783908071080791599664170866538088061836967103936433947787070068672183533326523316370911123438135644904401384304368230410467254086627981778866923497656678192846128933906906005269899725472218644305598944197162711932512512586839283057790743118806781283661695656893442275212703605794902759082850832855422405264957317988775866339476457721783729344116256656448187553012226037901478432392407066258367750572609107726924267077401439053067257842692103533450486770027223271173168400683811170698517540386225045073462986631764197259501842522014143479588918537751430535296446704645587831996792546724161485905310160385192649356313842727849372635722637482199850137101689284517030923518343256740250000301342390167883842891834006765996336031808405297299707857716795663044695031992127208249611812535012011793902720135085485308659443655312531205151690014075673266309273877194555286067807067691104638449253255021738338293454625291577650265728100234081790946412740602162577317063043073927182029634815614659862876926387500843068646634262299171729854596133118145726943409623601344641432949147576139889166633321082269125079085122835766104357949374044773003426572280289751827482028000836873049426022818821241988277563194748894176383694117448637877735136221689623671788382834137429652875033719989688861951297075973404809317272212203972745058092296971864444887728283022168460809262425094877877947981331032004152087673438429839969443648046171511976391809472889046554294619449443256544837882730750609620140287703706266956732614299826880081945638278754483208885197446770252694387520669319185740216584979686786479284471592166371135788852044938067651512340201884561114182111009221721583711634823053931209836287698742011019569469314604899124710413430348380972896453180679122594589193966571554368412024188727280953845410214481978047139552911011389399475665101116584868104041871908209567379691035055227105477864214172560544532490241457322472619207500646320636896015996165069321054361562787381277769133594506320403918394627551334173347077326993569012871149642387511646293380343836803581998178494982023934208479974329723042816908361595416421639496113407116195376728615338344430845543415433891932202110239174217200750428985923694160993291696311496545321898553424525033671722967134854525489557305503383109730508950223016481577727852892758358726821238098583316499526071246116630215078624704318877747124032976647421848424616365772031024921780184118017937939470856368254205085665455419393779570616990480077899378953525357650441598778031770509224081484432722073562288306817734286345366041658018323702592913498418994929783592094843266787158788954542273480182599806549652516322020213230250531787981976304684258638606006602529785807584045757102776578915172665344025150206366074306070134860653091115202904479304049108571728599967328808726993043869893238402860770705036237111579758398218501065238781228484793598701718459468441000510773972460337563318715351139669404349583445975346824017360189458283858100764540120478992987376082056633239328519362175098268067004925626009241774698523955460901002950589003735086726859835652391173486432590812241658574833408432030334138888909142612442552007742440812433817494156018521117106163863528047124588632765632108281146457624411967849782508779656993083097909976021312604185796108301737246840216283254050953956587078802131868836265366518748449614317340624246625405694568538247297950350842694333466802420526394419241078483203955084234221926600966240226174831609124067423564715321395005586731102579393800888638124912482092930320176522835705013376342355177091481815881136149063468926365664804864", 10)
				if z.Cmp(badNum) == 0 {
					continue
				}
				if tokens[tmp.index].Address == common.HexToAddress("0x34E89740adF97C3A9D3f63Cc2cE4a914382c230b") {
					log.Warn(z)
				}
				decimal := big.NewInt(10)
				decimal.Exp(decimal, big.NewInt(int64(tokens[tmp.index].Decimals)), nil)
				t := new(big.Float).SetInt(z)
				t = t.Quo(t, new(big.Float).SetInt(decimal))
				tokens[tmp.index].Balance = *t
				//tokens[tmp.index].BalanceStr = t.String()
				//log.Infoln(tokens[tmp.index].Address, tokens[tmp.index].BalanceStr, tmp.callRes.ReturnData)
				//if balance == float64
				//ten.Ex
				//t, ok := t.SetString(string(tmp.callRes.ReturnData))
				//if !ok {
				//	panic(tmp.callRes.ReturnData)
				//}
				//t.Quo(t,ten)
				//t, err := strconv.ParseFloat(string(tmp.callRes.ReturnData), 64)

				//if err != nil {
				//
				//	log.Fatal(err)
				//	tokens[tmp.index].Balance = 0.0
				//} else {
				//	tokens[tmp.index].Balance = t / float64(10^tokens[tmp.index].Decimals)
				//}

			}
		}

	}
	sort.Slice(tokens, func(i, j int) bool {
		//balance := tokens[i].Balance
		//tokens[j].Balance
		return tokens[i].Balance.Cmp(&tokens[j].Balance) > 0
	})
	ZERO := big.NewFloat(float64(0))
	firstZeroIndex := 0
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Balance.Cmp(ZERO) <= 0 {
			firstZeroIndex = i
			break
		}
	}
	return tokens[:firstZeroIndex]
	//res := make([]big.Int, len(tokens))

	//
	//token
	//for i := 0; i < len(chunkedCalls); i++ {
	//	//
	//	for j, _res := range (<-chunkChannel).callRes {
	//		z := new(big.Int)
	//		z.SetBytes(_res.ReturnData)
	//		res[i+j].Balance= *z
	//	}
	//}

}