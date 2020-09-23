#!/bin/bash

set -e

function yes_no_continue() {
    read -p "Are you sure? " -n 1 -r
    echo    # (optional) move to a new line
    if [[ ! $REPLY =~ ^[Yy]$ ]]
    then
        exit 1
    fi
}

# Ref: https://docs.docker.com/engine/install/ubuntu/
## Install docker-ce

function remove_old_version_docker_ce() {
    for pkg in `dpkg -l | grep -i docker | awk '{print $2}'`; do
        sudo apt-get remove -y ${pkg}
    done
    sudo apt-get remove -y docker containerd runc
}
remove_old_version_docker_ce

function install_new_version_docker_ce() {
    sudo apt-get update -y

    sudo apt-get install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg-agent \
        software-properties-common
    
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo apt-key fingerprint 0EBFCD88
    # 如果是在Linux Mint上安装docker, 需要手动替换 $(lsb_release -cs) 结果
    # 具体参考 https://www.linuxmint.com/download_all.php
    # 然后手动添加到 /etc/apt/sources.list.d/additional-repositories.list
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
                             
    sudo apt-get update -y

    sudo apt-cache policy docker-ce
    sudo apt-get install -y \
        docker-ce \
        docker-ce-cli \
        containerd.io

    sudo systemctl start docker
    sudo systemctl enable docker

    id -nG
    sudo usermod -aG docker ${USER}

    docker -v
    docker image ls
}
install_new_version_docker_ce

## Install docker-compose

function install_docker_compose() {
    sudo apt-get remove -y docker-compose

    DOCKER_COMPOSE_RELEASE=`curl -s https://github.com/docker/compose/releases/latest | cut -d'"' -f2 | cut -d'/' -f8-`
    sudo curl -L https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_RELEASE}/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    docker-compose -v
}
install_docker_compose

# Ref: https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#docker

## Install nvidia-docker

function install_nvidia_docker() {
    distribution=$(. /etc/os-release; echo $ID$VERSION_ID)
    curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
    # 如果是在Linux Mint上安装nvidia-docker, 需要手动替换 $distribution 结果
    # 具体参考 https://nvidia.github.io/nvidia-docker/
    curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list

    sudo apt-get update -y

    sudo apt-get install -y nvidia-docker2
    sudo systemctl restart docker
    sudo docker run --rm --gpus all nvidia/cuda:11.0-base nvidia-smi
}
install_nvidia_docker

sudo apt autoremove -y
