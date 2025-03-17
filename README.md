# dumpstat
Dump useful memory stats for every running process:

- Current memory use (**R**esident **S**et **S**ize)
- Max memory used (**H**igh **W**ater **M**ark of **R**esident **S**et **S**ize)
- Current swap usage

The measurements are in bytes.

## How to use it:

The help text is as such:

```
Usage of dumpstat:
  -c    csv format
  -i    influx line protocol
  -s    spaced format
  -ya
        yaml array
  -yh
        yaml hash using PID as keys
```

Spaced is handy for use with `column`:

```
$ dumpstat | grep -P 'brave|PID' | column -t
PID           NAME             STATE  VMRSS   VMHWM    VMSWAP
334883        brave            S      94260   148264   22700
336971        brave            S      79792   190348   49324
336961        brave            S      85136   158472   25132
927247        brave            S      29312   73472    9856
955824        brave            S      225596  699668   29840
948716        brave            S      385332  1478276  16764
```

You can feed it to telegraf:

```
[[exec]]
command = ["/usr/local/bin/dumpstat", "-i"]
data_format = "influx"
```

Output looks like this:
```
stat,name="brave",state=S rss=178532i,rss_hwm=240088i,swap=52472i 1742187345
stat,name="brave",state=S rss=45032i,rss_hwm=72644i,swap=11144i 1742187345
stat,name="brave",state=S rss=88128i,rss_hwm=123232i,swap=16188i 1742187345
stat,name="brave",state=S rss=109808i,rss_hwm=187300i,swap=18300i 1742187345
stat,name="brave",state=S rss=219184i,rss_hwm=699668i,swap=29840i 1742187345
```

You can get two types of YAML output too, array:

```
$ dumpstat -ya
---
- { pid: 927246, name: "brave", state: S, rss: 25728, rss_hwm: 73216, swap: 13184 }
- { pid: 948716, name: "brave", state: S, rss: 385572, rss_hwm: 1478276, swap: 16764 }
- { pid: 955824, name: "brave", state: S, rss: 219188, rss_hwm: 699668, swap: 29840 }
- { pid: 927290, name: "brave", state: S, rss: 45032, rss_hwm: 72644, swap: 11144 }
- { pid: 927279, name: "brave", state: S, rss: 109748, rss_hwm: 187300, swap: 18300 }
```

Or hash:
```
$ dumpstat -yh
---
936666: {name: "brave", state: S, rss: 42708, rss_hwm: 92800, swap: 13440 }
927228: {name: "brave-browser-s", state: S, rss: 3072, rss_hwm: 3584, swap: 256 }
948716: {name: "brave", state: S, rss: 385700, rss_hwm: 1478276, swap: 16764 }
955824: {name: "brave", state: S, rss: 216260, rss_hwm: 699668, swap: 29840 }
```
