# Node Exporter

This is an experiment as to how much meaningful kernel information can you collect from the proc fs.

Required reading: `man 5 proc`

Basic node exported that reads and exports some metrics from `/proc`

Files read:
* `/proc/diskstats` - disk IO metrics per volume
* `/proc/loadavg` - shows the load average for the whole machine
* `/proc/meminfo` - memory information
* `/proc/net/dev` - NIC interface level statistics - same as ifconfig
* `/proc/stats` - extracts CPU utilization and interupts
* `/proc/uptime` - reads uptime data
* `/proc/version` - kernel version string


## Notes
Tested manually with: `Linux version 6.6.26`
