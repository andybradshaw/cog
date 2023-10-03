import {DashboardBuilder} from "../../../generated/dashboard/dashboard/builder_gen";
import {TimePickerBuilder} from "../../../generated/dashboard/timepicker/builder_gen";
import {DashboardCursorSync, DashboardLinkType} from "../../../generated/types/dashboard/types_gen";


const builder = new DashboardBuilder("Some title")
    .uid("test-dashboard-codegen")
    .description("Some description")
    .time({from: "now-3h", to: "now"})
    .timepicker(
        new TimePickerBuilder()
            .refresh_intervals(["30s", "1m", "5m"])
    )
    .refresh("1m")
    .style("dark")
    .timezone("utc")
    .tooltip(DashboardCursorSync.Crosshair)
    .tags(["generated", "from", "cue"])
    .links([
        {
            // TODO: this is painful.
            title: "Some link",
            url: "http://google.com",
            type: DashboardLinkType.Link,
            tags: [],
            icon: "cloud",
            tooltip: "",
            asDropdown: false,
            targetBlank: false,
            includeVars: false,
            keepTime: false,
        },
    ]);

console.log(builder.build());
