package main

import (
	"fmt"
	"github.com/ctfang/command"
	"io/ioutil"
	"os"
	"strconv"
)
//扇区
type Sector struct {
	sectorBegin int//扇区开始
	sectorEnd int//扇区结束
}

func NewSector(sb int,se int) *Sector{
	return &Sector{
		sectorBegin: sb,
		sectorEnd:   se,
	}
}
type Track struct {
	data *[]byte //磁盘数据
	sectors map[int]Sector //扇区
}
//从指定扇区取出 扇区内容
func(tr *Track) PrintSector(num int) []byte{
	ses := tr.sectors[num]
	sd := (*tr.data)[ses.sectorBegin:ses.sectorEnd]
	return sd
}
//将文件写入扇区
func(tr *Track) setSector(num int,content []byte){
	ss := tr.sectors[num]
	for i,v := range content{
		(*tr.data)[ss.sectorBegin+i] =v
	}
}

func NewTrack(ss map[int]Sector,data []byte) Track{
	return Track{
		data:   &data,
		sectors: ss,
	}
}
//将内容写进 文件
func(tr *Track) WriteToSector(writecontent []byte,num int){
	var wrl int
	if len(writecontent) > 512 {
		wrl = len(writecontent) / 512
		for i := 0;i< wrl;i++{
			content := writecontent[(num + i ) * 512 :(num + i + 1 ) * 512]
			fmt.Println("文件大小:",len(content))
			if len(content) != 512{
				panic("错误的大小!")
			}
			tr.setSector(num + i,content)
		}
	}else{
		tr.setSector(num,writecontent)
	}
}
//将文件写出
func(tr *Track) WriteToImg(path string){
	WriteFile(*tr.data,path)
}
func main() {
	app := command.New()
	app.AddCommand(Echo{})
	app.Run()
}
// Echo 需要实现接口 CommandInterface
type Echo struct {
}

func (Echo) Configure() command.Configure {
	return command.Configure{
		Name:        "vhd",
		Description: "示例命令 hello",
		Input: command.Argument{
			// Argument参数为必须的输入的，不输入不执行
			// 匹配字符参数，匹配不到就是 value = false
			Argument: []command.ArgParam{
				{Name: "vhdfile", Description: "vhd文件路径~ 必填"},
			},
			Has: []command.ArgParam{
				{Name: "view", Description: "是否拥有 view 字符串!"},
				{Name: "vaild", Description: "引导分区自动添加识别数!"},
			},
			// 可选的参数，不输入也能执行
			Option:   []command.ArgParam{
				{Name: "w", Description: "要写入的文件!"},
				{Name: "o", Description: "要写出的文件!"},
				{Name: "n", Description: "要写入的扇区号!"},
			},
		},
	}
}

func (Echo) Execute(input command.Input) {

	//fmt.Println("hello")
	//fmt.Println("写入的文件",input.GetArgument("vhdfile"))
	//fmt.Println("操作的索引",input.GetOption("n"))
	//fmt.Println("要写入的扇区",input.GetOption("w"))
	//fmt.Println("是否输入了 view ：",input.GetHas("view"))
	datapath := input.GetOption("w")
	isview := input.GetHas("view")

	moshu := input.GetHas("vaild")

	path := input.GetArgument("vhdfile")
	output := input.GetOption("o")
	if path == ""{
		fmt.Println("file null Error")
	}
	index :=input.GetOption("n")

	vhdcontent := ReadFile(path)
	secotrs := make(map[int]Sector,0)
	cb := len(vhdcontent) / 512	//计算有多少个扇区
	for i:=1;i<=cb;i++{
		start := (i-1) * 512
		end := i * 512
		se := NewSector(start,end)
		secotrs[i-1]= *se
	}
	a := NewTrack(secotrs,vhdcontent)
	fmt.Println("磁盘扇区数:",len(a.sectors))
	fmt.Printf("磁盘容量: %dMB\n",len(a.sectors)/1024/2)
	//如果是观察模式
	if isview == true{
		if index == ""{
			fmt.Println("please input index of vhd to view!")
			return
		}
		nt, _ := strconv.Atoi(index)
		data :=a.PrintSector(nt)
		fmt.Printf("start byte:%d end byte:%d\n",nt * 512,nt * 512 + 512)
		fmt.Printf("    O(∩_∩)O       +2   +4   +6   +8  +10   +12 +14  +16")
		for i,m := range data{
			if i % 16 == 0{
				fmt.Println()
				fmt.Printf("OffSet:%5d ==>",i)
			}
			if i % 2 ==0{
				fmt.Print(" ")
			}
			fmt.Printf("%02x",m)
		}
		fmt.Println()
	}else if  datapath != "" {
		if index ==""{
			fmt.Println("please input index of vhd to Write!")
			return
		}
		nt, _ := strconv.Atoi(index)
		writedata := ReadFile(datapath)
		fmt.Println("Data Will Write To VHD FIle.",writedata)
		//55 AA 为intelCpu 识别引导区的魔数  如果不加如法加载
		if moshu && nt == 0{
			for len(writedata) != 512{
				writedata= append(writedata, 0x00)
			}
			writedata[510]= 0x55
			writedata[511]= 0xAA
		}
		a.WriteToSector(writedata,nt)
		//写出文件 如果没指定 就覆盖源文件
		if output != ""{
			a.WriteToImg(output)
		}else{
			a.WriteToImg(path)
		}
	}

}

func ReadFile(path string)  ([]byte){
	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("read fail", err)
	}
	return f
}

func WriteFile(content []byte,path string){
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	n2, err := f.Write(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(n2, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}