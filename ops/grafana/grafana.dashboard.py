from grafanalib.core import *

dashboard = Dashboard(
    title="GoToURL Stats",
    rows=[
        Row(panels=[
            Graph(
                title="GoToURL QPS",
                dataSource='k8s-prometheus',
                targets=[
                    Target(
                        expr='sum(rate(http_request_duration_seconds_count{job="gotourl-app",code=~"1.."}[1m]))',
                        legendFormat="1xx",
                        refId='A',
                    ),
                    Target(
                        expr='sum(rate(http_request_duration_seconds_count{job="gotourl-app",code=~"2.."}[1m]))',
                        legendFormat="2xx",
                        refId='B',
                    ),
                    Target(
                        expr='sum(rate(http_request_duration_seconds_count{job="gotourl-app",code=~"3.."}[1m]))',
                        legendFormat="3xx",
                        refId='C',
                    ),
                    Target(
                        expr='sum(rate(http_request_duration_seconds_count{job="gotourl-app",code=~"4.."}[1m]))',
                        legendFormat="4xx",
                        refId='D',
                    ),
                    Target(
                        expr='sum(rate(http_request_duration_seconds_count{job="gotourl-app",code=~"5.."}[1m]))',
                        legendFormat="5xx",
                        refId='E',
                    ),
                ],
                yAxes=[
                    YAxis(format=OPS_FORMAT),
                    YAxis(format=SHORT_FORMAT),
                ],
                alert=Alert(
                    name="Too many 500s",
                    message="More than 5 QPS of HTTP 5xx for 5 minutes",
                    frequency="60s",
                    noDataState=ALERTLIST_STATE_OK,
                    alertConditions=[
                        AlertCondition(
                            Target(
                                legendFormat="5xx",
                                refId='E'  # 5xx metric
                            ),
                            timeRange=TimeRange("5m", "now"),
                            evaluator=GreaterThan(5),
                            operator=OP_AND,
                            reducerType=RTYPE_COUNT,
                        ),
                    ]
                )
                # todo nsokil grafana doesn't support separate alerts in 1 graph
                #  https://github.com/grafana/grafana/issues/7832
                # ,
                # Alert(
                #     name="High Request Rate",
                #     message="More than 1000 QPS for 5 minutes",
                #     alertConditions=[
                #         AlertCondition(
                #             Target(
                #                 expr='sum(rate(http_request_duration_seconds_count{job="gotourl-app"}[1m]))',
                #                 legendFormat="5xx",
                #                 refId='A',
                #             ),
                #             timeRange=TimeRange("5m", "now"),
                #             evaluator=GreaterThan(1000),
                #             operator=OP_AND,
                #             reducerType=RTYPE_SUM,
                #         ),
                #     ])]
            ),
            Graph(
                title="GoToURL latency",
                dataSource='k8s-prometheus',
                targets=[
                    # todo nsokil the ">= 0" part is due to PromQL returning NaN, check newer Grafana/Prometheus versions
                    #  https://github.com/grafana/grafana/issues/11512
                    Target(
                        expr='histogram_quantile(0.5, sum(rate(http_request_duration_seconds_bucket{job="gotourl-app"}[1m])) by (le)) >= 0',
                        legendFormat="0.5 quantile",
                        refId='A',
                    ),
                    Target(
                        expr='histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{job="gotourl-app"}[1m])) by (le)) >= 0',
                        legendFormat="0.95 quantile",
                        refId='B',
                    ),
                    Target(
                        expr='histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket{job="gotourl-app"}[1m])) by (le)) >= 0',
                        legendFormat="0.99 quantile",
                        refId='C',
                    ),
                ],
                yAxes=single_y_axis(format=SECONDS_FORMAT),
                alert=Alert(
                    name="High p99 latency",
                    message="p99 is above 100ms for 5 minutes",
                    noDataState=ALERTLIST_STATE_OK,
                    alertConditions=[
                        AlertCondition(
                            Target(
                                legendFormat="0.99 quantile",
                                refId='C'  # 0.99 metric
                            ),
                            timeRange=TimeRange("5m", "now"),
                            evaluator=GreaterThan(0.1),
                            operator=OP_AND,
                            reducerType=RTYPE_AVG,
                        ),
                    ],
                )
            ),
        ]),
    ],
).auto_panel_ids()
