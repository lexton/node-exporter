# Node Exporter

Basic node exported that reads and exports some metrics from `/proc`

Files read:
* `/proc/stats` - extracts CPU and interupts
* `/proc/uptime` - reads uptime data
* `/proc/diskstats` - disk IO metrics per volume
* `/proc/loadavg` - shows the load average for the whole machine
* `/proc/version` - renders the kernel version
* `/proc/net/dev` - renders NIC interface level statistics

# Details

For the exact spec see:

`man 5 proc` 

Tested with: `Linux version 6.6.26`
