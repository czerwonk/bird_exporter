## BIRD Grafana Dashboard

## BIRD Grafana Dashboard

There are three dashboards for this exporter. One is a more general RS dashboard, and two more advanced dashboards that focus primarily on **BGP** and **OSPF & BFD** protocols.

All three of them are available in the dashboards folder. Just copy the JSON file and import it into your Grafana instance.

Some notes about dashboards:

1) The RS dashboard is also available in the [Grafana cloud](https://grafana.com/grafana/dashboards/5259-bird-rs) as well.


<image src="./img/bird_exporter.png"></image>

2) For the two **BGP** and **OSPF & BFD** dashboards:
  - The default port is `9324`, so if you have changed that, please don't forget to change it as well from metrics (the `instance` label).

Why does the **OSPF & BFD** dashboard contain both of them together?
The BFD is a fast failure detection protocol that can be used in conjunction with OSPF to provide faster failure detection than OSPF Hello messages alone. It is common in most OSPF setups to use BFD for failure detection. That's why it has been included there as well.

The BGP dashboard:

<image src="./img/bgp.png"></image>

The OSPF & BFD dashboard:

<image src="./img/ospf.png"></image>
