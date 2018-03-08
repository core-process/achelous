# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/artful64"
  config.vm.provision "file", source: ".build/dist/achelous_1.0-1_amd64.deb", destination: "achelous_1.0-1_amd64.deb"
  config.vm.provision "shell", inline: <<-SHELL
    dpkg -i achelous_1.0-1_amd64.deb
  SHELL
end
