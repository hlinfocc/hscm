# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

hscm_version=`./bin/hscm -v`
echo "build version: ${hscm_version}"

if [ -d "./dist" ];then
    rm -rf ./dist
fi

mkdir -p ./dist/packages

# cross_compiles
make -f ./Makefile.cross-compiles


os_all='linux'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64'

cd ./dist

for os in $os_all; do
    for arch in $arch_all; do
        echo "package: OS:${os},arch: ${arch}"
        hscm_dir_name="hscm_${hscm_version}_${os}_${arch}"
        hscm_path="./packages/hscm_${hscm_version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ];then
            if [ ! -f "./hscm_${os}_${arch}.exe" ];then
                continue
            fi
            mkdir -p ${hscm_path}
            mv ./hscm_${os}_${arch}.exe ${hscm_path}/hscm.exe
        else
             if [ ! -f "./hscm_${os}_${arch}" ];then
                continue
            fi
            mkdir -p ${hscm_path}
            mv ./hscm_${os}_${arch} ${hscm_path}/hscm
        fi

        cp ../LICENSE ${hscm_path}
        if [ "x${os}" = x"linux" ]; then
            \cp ../conf/linux/* ${hscm_path}/
        fi

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${hscm_dir_name}.zip ${hscm_dir_name}
        else
            tar -zcvf ${hscm_dir_name}.tar.gz ${hscm_dir_name}
        fi
        cd ..
        rm -rf ${hscm_path}
    done
done

\cp ./packages/* ./

if [ -d "./packages" ];then
    \rm -rf ./packages
fi

