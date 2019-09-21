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
                        expr='sum(rate(promhttp_metric_handler_requests_total{job="gotourl-app",code=~"1.."}[1m]))',
                        legendFormat="1xx",
                        refId='A',
                    ),
                    Target(
                        expr='sum(rate(promhttp_metric_handler_requests_total{job="gotourl-app",code=~"2.."}[1m]))',
                        legendFormat="2xx",
                        refId='B',
                    ),
                    Target(
                        expr='sum(rate(promhttp_metric_handler_requests_total{job="gotourl-app",code=~"3.."}[1m]))',
                        legendFormat="3xx",
                        refId='C',
                    ),
                    Target(
                        expr='sum(rate(promhttp_metric_handler_requests_total{job="gotourl-app",code=~"4.."}[1m]))',
                        legendFormat="4xx",
                        refId='D',
                    ),
                    Target(
                        expr='sum(rate(promhttp_metric_handler_requests_total{job="gotourl-app",code=~"5.."}[1m]))',
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
                    message="More than 5 QPS of 500s for 5 minutes",
                    alertConditions=[
                        AlertCondition(
                            Target(
                                expr='sum(rate(promhttp_metric_handler_requests_total{job="gotourl-app",code=~"5.."}[1m]))',
                                legendFormat="5xx",
                                refId='A',
                            ),
                            timeRange=TimeRange("5m", "now"),
                            evaluator=GreaterThan(5),
                            operator=OP_AND,
                            reducerType=RTYPE_SUM,
                        ),
                    ],
                )
            ),
            Graph(
                title="GoToURL latency",
                dataSource='k8s-prometheus',
                targets=[
                    Target(
                        expr='histogram_quantile(0.5, sum(rate(http_request_duration_seconds_bucket{job="gotourl-app"}[5m])) by (le))',
                        legendFormat="0.5 quantile",
                        refId='A',
                    ),
                    Target(
                        expr='histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket{job="gotourl-app"}[5m])) by (le))',
                        legendFormat="0.99 quantile",
                        refId='B',
                    ),
                ],
                yAxes=single_y_axis(format=SECONDS_FORMAT),
            ),
        ]),
    ],
).auto_panel_ids()
