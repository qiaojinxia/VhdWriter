A sample tool writer by go,your can use it to Writer Img with VHD file system。

useage:

>./main vhd caomao.vhd -n=0 view

Vhd is default parameters no need to cahnge。

cammao.vhd is necessary parameters be used to writing disk。

-n represent index  where sector can Writer or Read。

view represent model Read

The command above can show up the vhd file range 512 bytes  every time。 offfset by n ~ (n+1)*512

<img src="/Users/qiao/Library/Application Support/typora-user-images/image-20200126194956811.png" alt="image-20200126194956811" style="zoom:40%;" />

> ./ads vhd caomao.vhd -n=0 -w=/Users/qiao/VirtualBox\ VMs/cbxos/a.bin

-w parameters represent data to be written vhd file。

-o If you specify this parameter, a new vhdfile will be written out in the specified file path if you do not specify to overwrite the original file。

<img src="/Users/qiao/Library/Application Support/typora-user-images/image-20200126200828738.png" alt="image-20200126200828738" style="zoom:50%;" />

