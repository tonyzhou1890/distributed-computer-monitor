// @ts-nocheck
const $ = document.querySelector.bind(document)
let host = null, // 当前主机
    hostList = [], // 主机列表
    timeChecked = 1, // 当前时间范围
    defaultDataMap = {
      xAxis: [],
      cpu: {
        model: '',
        core: 0,
        load: [],
        temp: []
      },
      mem: {
        cap: 0,
        load: []
      },
      swap: {
        cap: 0,
        load: []
      }
    },
    dataMap = JSON.parse(JSON.stringify(defaultDataMap))

window.onload = () => {
  getHostList()
}

// 获取主机列表
function getHostList() {
  fetch('./api/web/host/list')
    .then(res => {
      res.json()
        .then(data => {
          if (data.code === 0) {
            hostList = data.data ?? []
            if (hostList.length) {
              host = hostList[0].id
            }
            // 渲染主机列表
            renderHostList()
            // 设置时间
            setTime()
          } else {
            alert(data.msg)
          }
        })
    })
    .catch(e => {
      alert(e.message)
    })
}

// 渲染主机列表
function renderHostList() {
  const hostContainer = $('#hostList')
  if (!hostContainer) {
    alert('无法渲染主机列表')
    return
  }
  const htmlstr = template('selectTpl', { hostList })
  hostContainer.innerHTML = htmlstr
  const hostSelect = $('#hostSelect')
  if (hostSelect) {
    hostSelect.value = host
    hostSelect.onchange = handleHostChange
  }
  getHostData()
}

// 主机切换
function handleHostChange(ev) {
  console.log(ev.target.value)
  host = Number(ev.target.value)
  getHostData()
}

// 设置时间
function setTime() {
  const timeSelect = $('#timeSelect')
  if (timeSelect) {
    timeSelect.setAttribute('data-time', timeChecked)
  }
}

// 时间改变
function handleTimeChange(value) {
  timeChecked = value
  setTime()
  getHostData()
}

// 获取主机信息
function getHostData() {
  let startTime = dayjs()
  switch (timeChecked) {
    case 1:
      startTime = startTime.subtract(1, 'hour')
      break
    case 2:
      startTime = startTime.subtract(1, 'day')
      break
    case 3:
      startTime = startTime.subtract(1, 'month')
      break
    case 4:
      startTime = startTime.subtract(3, 'month')
      break
    default:
      break
  }
  // console.log(startTime, timeChecked)
  fetch(`./api/web/record/list?hostId=${host}&startTime=${startTime.format('YYYY-MM-DD HH:mm:ss')}&endTime=${dayjs().format('YYYY-MM-DD HH:mm:ss')}`, {
    method: 'GET'
  })
    .then(res => {
      res.json()
        .then(data => {
          if (data.code === 0) {
            const list = data.data || []
            dataMap = JSON.parse(JSON.stringify(defaultDataMap))
            // 信息处理
            list.map(item => {
              dataMap.xAxis.push(dayjs(item.created).format('MM:DD HH:mm'))
              dataMap.cpu.model = item.cpuType
              dataMap.cpu.core = item.cpuCore
              dataMap.cpu.load.push(toAccuracy(item.cpuLoad, 1))
              dataMap.cpu.temp.push(toAccuracy(item.cpuTemp, 1))
              dataMap.mem.cap = toAccuracy(item.memCap / 1024 / 1024, 3)
              dataMap.mem.load.push(toAccuracy(item.memLoad, 1))
              dataMap.swap.cap = toAccuracy(item.swapCap / 1024 / 1024, 3)
              dataMap.swap.load.push(toAccuracy(item.swapLoad, 1))
            })
            // 显示相关信息
            showBaseInfo()
            // 绘制图表
            renderChart()
          } else {
            alert(data.msg)
          }
          // console.log(data)
        })
    })
    .catch(e => {
      alert(e.message)
    })
}

// 显示基础信息
function showBaseInfo() {
  // cpu
  $('#cpuModel').innerText=`型号：${dataMap.cpu.model ?? ''}`
  $('#cpuCore').innerText = `核心：${dataMap.cpu.core ?? 0}`
  // 内存
  $('#memCap').innerText = `容量：${dataMap.mem.cap}GB`
  // swap
  $('#swapCap').innerText = `容量：${dataMap.swap.cap}GB`
}

// 绘制图表
function renderChart() {
  // 基础配置
  const baseConfig = {
    title: {
      text: ''
    },
    tooltip: {
      trigger: 'axis',
    },
    legend: {},
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dataMap.xAxis
    },
    yAxis: {
      type: 'value'
    },
    dataZoom: [
      {
      }
    ],
    series: [
      {
        data: [],
      name: '',
        type: 'line'
      }
    ]
  }
  // cpu 负载
  const cpuLoadConfig = JSON.parse(JSON.stringify(baseConfig))
  cpuLoadConfig.series[0].data = dataMap.cpu.load
  cpuLoadConfig.title.text = 'cpu 负载'
  cpuLoadConfig.yAxis.name = '%'
  if (!window.cpuLoadChartIns) {
    window.cpuLoadChartIns = echarts.init($('#cpuLoadChart'))
  }
  window.cpuLoadChartIns.setOption(cpuLoadConfig)
  
  // cpu 温度
  const cpuTempConfig = JSON.parse(JSON.stringify(baseConfig))
  cpuTempConfig.series[0].data = dataMap.cpu.temp
  cpuTempConfig.title.text = 'cpu 温度'
  cpuTempConfig.yAxis.name = '℃'
  if (!window.cpuTempChartIns) {
    window.cpuTempChartIns = echarts.init($('#cpuTempChart'))
  }
  window.cpuTempChartIns.setOption(cpuTempConfig)
  
  // 内存负载
  const memCapConfig = JSON.parse(JSON.stringify(baseConfig))
  memCapConfig.series[0].data = dataMap.mem.load
  memCapConfig.title.text = '内存负载'
  memCapConfig.yAxis.name = '%'
  if (!window.memLoadChartIns) {
    window.memLoadChartIns = echarts.init($('#memLoadChart'))
  }
  window.memLoadChartIns.setOption(memCapConfig)

  // swap 负载
  const swapCapConfig = JSON.parse(JSON.stringify(baseConfig))
  swapCapConfig.series[0].data = dataMap.swap.load
  swapCapConfig.title.text = 'swap 负载'
  swapCapConfig.yAxis.name = '%'
  if (!window.swapLoadChartIns) {
    window.swapLoadChartIns = echarts.init($('#swapLoadChart'))
  }
  window.swapLoadChartIns.setOption(swapCapConfig)
}

// 小数截取
function toAccuracy(val, n) {
  val = typeof val === 'number' ? val : 0
  n = typeof n === 'number' ? n : 0
  return Number(val.toFixed(n))
}