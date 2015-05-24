# -*- mode: ruby -*-
# vi: set ft=ruby :

def set_hostname(host, domain)
  set_hostname = <<SCRIPT
    sudo hostname #{host}.#{domain}
SCRIPT
end

def fake_dns(num_servers, num_nodes, domain)
  script = ""
  num_servers.each do |i|
    script = script + "echo \"192.168.33.1#{i}  consul-#{i}.#{domain}\" >> /etc/hosts\n"
  end
  num_nodes.each do |i|
    script = script + "echo \"192.168.33.2#{i}  node-#{i}.#{domain}\" >> /etc/hosts\n"
  end
  return script
end

def consul_conf(type, domain, num = 0)
  case type
  when "server"
    c = <<SCRIPT
{
  "bootstrap": true,
  "server": true,
  "data_dir": "/var/lib/consul",
  "ui_dir": "/opt/consului",
  "log_level": "info",
  "datacenter": "local",
  "bind_addr": "192.168.33.1#{num}",
  "client_addr": "192.168.33.1#{num}",
  "addresses": {
    "dns": "192.168.33.1#{num}",
    "http": "192.168.33.1#{num}",
    "rpc":"192.168.33.1#{num}"
  }
}
SCRIPT
  when "client"
    c = <<SCRIPT
{
  "server": false,
  "data_dir": "/var/lib/consul",
  "log_level": "info",
  "datacenter": "local",
  "start_join": ["consul-1.#{domain}"]
}
SCRIPT
  when "upstart"
    c = <<SCRIPT
description "Consul agent"

start on runlevel [2345]
stop on runlevel [!2345]

respawn

script
  if [ -f "/etc/service/consul" ]; then
    . /etc/service/consul
  fi

  # Make sure to use all our CPUs, because Consul can block a scheduler thread
  export GOMAXPROCS=`nproc`

  # Get the public IP
  BIND=`ifconfig eth1 | grep "inet addr" | awk '{ print substr($2,6) }'`

  exec /usr/bin/consul agent \
    -config-dir="/etc/consul.d" \
    -bind=$BIND \
    ${CONSUL_FLAGS} \
    >>/var/log/consul.log 2>&1
end script
SCRIPT
  end
  return c
end

def install_consul(type, domain, num)
  install_consul = <<SCRIPT
echo Installing dependencies...
sudo apt-get install -y unzip curl
if [ -f /tmp/consul.zip ]; then
  echo Consul already installed ...
else
  echo Fetching Consul...
  cd /tmp/
  wget https://dl.bintray.com/mitchellh/consul/0.5.2_linux_amd64.zip -O consul.zip
  echo Installing Consul...
  unzip consul.zip
  sudo chmod +x consul
  sudo mv consul /usr/bin/consul
fi
echo Consul Configuration
sudo mkdir -p /etc/consul.d
sudo echo #{consul_conf(type, domain, num)} /etc/consul.d/config.json
sudo chown root:root /etc/consul.d/*
sudo chmod 644 /etc/consul.d/*
echo Consul upstart Installation
sudo echo #{consul_conf("upstart", domain)} /etc/init/consul.conf
sudo chown root:root /etc/init/consul.conf
echo Consul Agent Start
sudo service consul restart
SCRIPT

if type == "server"
install_ui = <<SCRIPT
if [ -f /tmp/consului.zip ]; then
  echo Consul UI already installed ...
else
  echo Fetching Consul UI ...
  cd /tmp/
  wget https://dl.bintray.com/mitchellh/consul/0.5.2_web_ui.zip -O consului.zip
  echo Installing Consul UI...
  unzip consului.zip
  sudo mv dist /opt/consului
fi
SCRIPT
  install_consul = install_consul + install_ui
end

  return install_consul
end

def motd(host, domain)

  motd = <<SCRIPT

sudo apt-get install -y figlet
figlet -f small fwrules > /etc/motd
echo >> /etc/motd
figlet -f small #{host} >> /etc/motd
echo "\n\n#{host}.#{domain}" >> /etc/motd

SCRIPT

end



Vagrant.configure(2) do |config|

  consul_servers = 1
  consul_nodes = 4

  domain = "fwrules.local"

  (1..consul_servers).each do |id|
    config.vm.define h="consul-#{id}" do |v|
      v.vm.box = "chef/debian-7.8"
      v.vm.network "private_network", ip: "192.168.33.#{10+id}"
      v.vm.network "forwarded_port", guest: 8500, host: 8500 + id
      v.vm.hostname = "#{h}.#{domain}"

      v.vm.provider "virtualbox" do |vb|
        vb.customize ["modifyvm", :id, "--memory", "256"]
        vb.cpus = 4
      end
      v.vm.provision "shell", inline: set_hostname(h, domain)
      v.vm.provision "shell", inline: fake_dns(consul_servers, consul_nodes, domain)
      v.vm.provision "shell", inline: install_consul("server", domain, id)
      v.vm.provision "shell", inline: motd(h, domain)

      v.vm.provision "shell", inline: $env_vars
      v.vm.provision "shell", inline: $install_go
    end
  end

  (1..consul_nodes).each do |id|
    config.vm.define h="node-#{id}" do |v|
      v.vm.box = "chef/debian-7.8"
      v.vm.network "private_network", ip: "192.168.33.#{20+id}"
      v.vm.network "forwarded_port", guest: 49999, host: 48000 - id
      v.vm.hostname = "#{h}.#{domain}"

      v.vm.provider "virtualbox" do |vb|
        vb.customize ["modifyvm", :id, "--memory", "128"]
        vb.cpus = 4
      end
      v.vm.provision "shell", inline: set_hostname(h, domain)
      v.vm.provision "shell", inline: fake_dns(consul_servers, consul_nodes, domain)
      v.vm.provision "shell", inline: install_consul("client", domain, id)
      v.vm.provision "shell", inline: motd(h, domain)

      v.vm.provision "shell", inline: set_hostname(h, domain)
      v.vm.provision "shell", inline: $fake_dns
      v.vm.provision "shell", inline: $env_vars
      v.vm.provision "shell", inline: install_consul("client")
      v.vm.provision "shell", inline: $install_go
      v.vm.provision "shell", inline: install_docker(h, domain)
      v.vm.provision "shell", inline: $install_sxnode
      v.vm.provision "shell", inline: pull_images(h, domain)
    end
  end

end