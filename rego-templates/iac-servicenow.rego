package postee.iac.servicenow

import data.postee.by_flag
import data.postee.with_default
import future.keywords
import future.keywords.if

################################################ Templates ################################################
# Template is used in `work notes`.
html_tpl:=`
<p><b>Repository Name:</b> %s</p>
<p> </p>
<!-- Stats -->
<h3> Vulnerability summary: </h3>
%s
<h3> Misconfiguration summary: </h3>
%s
<h3> Pipeline Misconfiguration summary: </h3>
%s
<p><b>Resourse policy name:</b> %s</p>
<p><b>Resourse policy application scopes:</b> %s</p>
`

summary_tpl =`Registry name: %s`

#Extra % is required in width:100%
table_tpl:=`
<TABLE border='1' style='width: 100%%; border-collapse: collapse;'>
%s
</TABLE>
`

cell_tpl:=`<TD style='padding: 5px;'>%s</TD>
`

row_tpl:=`
<TR>
%s
</TR>`

colored_text_tpl:="<span style='color:%s'>%s</span>"

############################################## Html rendering #############################################
render_table(content_array) = s {
	rows := [tr |
    			cells:=content_array[_]
    			tds:= [td |
                	ctext:=cells[_]
                    td := to_cell(ctext)
                ]
                tr=sprintf(row_tpl, [concat("", tds)])
    		]

	s:=sprintf(table_tpl, [concat("", rows)])
}

to_cell(txt) = c {
    c:= sprintf(cell_tpl, [txt])
}

to_colored_text(color, txt) = spn {
    spn :=sprintf(colored_text_tpl, [color, txt])
}

####################################### Template specific functions #######################################
to_severity_color(color, level) = spn {
 spn:=to_colored_text(color, format_int(with_default(input,level,0), 10))
}

severities_stats(critical, high, medium, low, unknown) := [
                        ["critical", to_severity_color("#c00000", critical)],
                        ["high", to_severity_color("#e0443d", high)],
                        ["medium", to_severity_color("#f79421", medium)],
                        ["low", to_severity_color("#e1c930", low)],
                        ["unknown", to_severity_color("green", unknown)]
                    ]

vulnerability_stats = stats{
    stats := severities_stats("vulnerability_critical_count",
                              "vulnerability_high_count",
                              "vulnerability_medium_count",
                              "vulnerability_low_count",
                              "vulnerability_unknown_count")
}

misconfiguration_stats = stats{
    stats := severities_stats("misconfiguration_critical_count",
                              "misconfiguration_high_count",
                              "misconfiguration_medium_count",
                              "misconfiguration_low_count",
                              "misconfiguration_unknown_count")
}

pipeline_stats = stats{
    stats := severities_stats("pipeline_misconfiguration_critical_count",
                              "pipeline_misconfiguration_high_count",
                              "pipeline_misconfiguration_medium_count",
                              "pipeline_misconfiguration_low_count",
                              "pipeline_misconfiguration_unknown_count")
}

############################################## result values #############################################
title = sprintf(`Aqua security | Repository | %s | Scan report`, [input.repository_name])

result_assigned_to := by_flag(input.application_scope_owners[0], "", count(input.application_scope_owners) == 1)
result_assigned_group := by_flag(input.application_scope[0], "", count(input.application_scope) == 1)

result_severity := 1 if {
    input.vulnerability_critical_count +
    input.misconfiguration_critical_count +
    input.pipeline_misconfiguration_critical_count > 0
} else = 2 if {
    input.vulnerability_high_count +
    input.misconfiguration_high_count +
    input.pipeline_misconfiguration_high_count > 0
} else = 3

result_summary := summary{
    summary := sprintf(summary_tpl, [input.repository_name])
}

result = msg {

    msg := sprintf(html_tpl, [
    input.repository_name,
    render_table(vulnerability_stats),
    render_table(misconfiguration_stats),
    render_table(pipeline_stats),
    with_default(input, "response_policy_name", "none"),
    with_default(input, "application_scope", "none")
    ])
}