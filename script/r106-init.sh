sh /home/root/r106-enhance/script/r106-enhance-monitor.sh > /dev/null 2>&1 &
#sh /home/root/frp/frp-monitor.sh > /dev/null 2>&1 &
#sleep 5
/home/root/frp/frpc -c /home/root/frp/frpc.toml > /dev/null 2>&1 &