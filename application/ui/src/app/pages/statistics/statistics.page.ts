import { Component, ElementRef, OnInit, ViewChild } from "@angular/core";
import Chart from "chart.js/auto";
import { formatDate } from "@angular/common";
import {
  MatchingStat,
  MatchingStatService,
} from "../../services/matching-stat.service";
import { TranslateService } from "@ngx-translate/core";

@Component({
  selector: "app-statistics",
  templateUrl: "./statistics.page.html",
  styleUrls: ["./statistics.page.scss"],
})
export class StatisticsPage implements OnInit {
  @ViewChild("chartCanvas") chartCanvas: ElementRef;
  data: any = [];
  canvasChart: Chart;

  matchingstatistics: MatchingStat[] = [];
  connectionsGenerated: string;
  conversationsWellTriggered: string;
  conversationsNotSent: string;
  numUsers: boolean = false;
  totalNumConnections = 0;
  numconves = 0;
  numfailed = 0;
  numMatchingSection = 0;

  constructor(
    private matchingstatService: MatchingStatService,
    public translate: TranslateService
  ) {}

  ngOnInit() {
    this.connectionsGenerated = this.translate.instant(
      "NUM_CONNECTIONS_GENERATED"
    );
    this.conversationsWellTriggered = this.translate.instant(
      "CONNECTIONS_INVITE_SENT"
    );
    this.conversationsNotSent = this.translate.instant(
      "CONNECTIONS_INVITE_NOT_SENT"
    );
    this.getMathingStats();
  }

  getMathingStats() {
    let numgroups = [];
    let numemails = [];
    let numfailed = [];
    let labels = [];

    this.matchingstatService.getMatchingStats().subscribe((matchingStats) => {
      if (matchingStats != null) {
        this.numMatchingSection = matchingStats.length;

        this.matchingstatistics = matchingStats;
        for (const mathingStat of matchingStats) {
          labels.push(
            formatDate(
              mathingStat.createdAt,
              "MMM d, y hh:mm:ss",
              "en" /**navigator.language.split("-")[0] ||  */
            )
          );
          numgroups.push(mathingStat.numGroups);
          numemails.push(mathingStat.numConversations);
          numfailed.push(mathingStat.numFailed);
          this.numconves += mathingStat.numConversations;
          this.numfailed += mathingStat.numFailed;
          this.totalNumConnections += mathingStat.numGroups;
        }
        this.data = {
          labels: labels,
          datasets: [
            {
              label: this.connectionsGenerated,
              data: numgroups,
              backgroundColor: "rgba(56,128,255,1)",
              borderWidth: 1,
            },
            {
              label: this.conversationsWellTriggered,
              data: numemails,
              backgroundColor: "rgba(45,211,111,1)",
              borderWidth: 1,
            },
            {
              label: this.conversationsNotSent,
              data: numfailed,
              backgroundColor: "rgba(255,196,9,1)",
              borderWidth: 1,
            },
          ],
        };
        this.changeChart({
          detail: {
            value: "bar",
          },
        });
      }
    });
  }

  changeChart(event: any) {
    const type = event.detail.value || "bar";
    if (this.canvasChart) {
      this.canvasChart.destroy();
    }
    this.canvasChart = new Chart(this.chartCanvas.nativeElement, {
      type: type,
      data: this.data,
      options: {
        maintainAspectRatio: false,
        indexAxis: "x",
        plugins: {
          legend: {
            align: "start",
            position: "top",
          },
        },
      },
    });
  }

  addRandom(points: any): number {
    return Number(points) - Number(Math.floor(Math.random() * 100 + 1));
  }
}
