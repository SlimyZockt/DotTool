#!/usr/bin/env bash


pushd $1

if [ $# -eq 0 ]
  then
    echo "No arguments supplied"
fi

rm -rf ./niri/
rm -rf ./hypr/
rm -f  ./tmux/tmux.conf
rm -rf ./waybar/
rm -rf ./wezterm/
rm -rf ./fish/
rm -rf ./ghostty/
rm -rf ./fuzzel/

cp -r ~/.config/niri/ .
cp -r ~/.config/hypr/ .
cp -r ~/.config/tmux/tmux.conf ./tmux/tmux.conf
cp -r ~/.config/waybar/ .
cp -r ~/.config/wezterm/ .
cp -r ~/.config/fish/ .
cp -r ~/.config/ghostty/ .
cp -r ~/.config/fuzzel/ .

git add .
git commit -am "$(($(git log -1 --pretty=%B)+1))"
git push
popd
