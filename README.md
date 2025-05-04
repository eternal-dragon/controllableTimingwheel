# 可调速时间轮 Timing Wheel with Adjustable Time Scale
## 项目概述 Project Overview
本项目实现了一个具有可调节时间流速的层级时间轮。它提供了一种高效管理不同时间尺度上定时事件的方法，并允许用户根据应用场景动态调整时间流速。

This project implements a hierarchical timing wheel with adjustable time flow rate. It provides an efficient way to manage timer events across different time scales and allows for dynamic adjustment of the time flow rate to adapt to various application scenarios.

## 核心功能 Key Features

* 层级时间轮结构 ：由多个不同分辨率的时间轮层组成。底层具有较细的时间粒度，而高层具有较粗的时间粒度，能够高效地管理较宽时间范围内的定时事件。
* Hierarchical Timing Wheel Structure : Consists of multiple layers of timing wheels with different resolutions. The lower layers have finer time granularity, while the higher layers have coarser time granularity, enabling efficient management of timer events across a wide time range.

* 可调节时间流速 ：提供接口以动态调整时间流速，使用户能够根据需求加快或减慢系统内的时间流逝速度。该功能在模拟、测试环境等需要时间控制的场景中尤为有用。
* Adjustable Time Scale : Provides an interface to dynamically adjust the time flow rate, allowing users to speed up or slow down the passage of time in the system according to their needs. This feature is particularly useful in simulations, testing environments, and other scenarios where time control is required.






## 安装 Installation

```bash
$ go get -u github.com/eternal-dragon/controllableTimingwheel
```

## 文档 Documentation

For usage and examples see the [Godoc][1].

##  License

[MIT][2]

[1]: https://godoc.org/github.com/eternal-dragon/controllableTimingwheel
[2]: http://opensource.org/licenses/MIT
