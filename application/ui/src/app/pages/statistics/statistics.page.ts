import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import Chart from 'chart.js/auto';
import { formatDate } from '@angular/common';
import { MatchingStat, MatchingStatService } from '../../services/matching-stat.service';
import { TranslateService } from '@ngx-translate/core';


@Component({
  selector: 'app-statistics',
  templateUrl: './statistics.page.html',
  styleUrls: ['./statistics.page.scss'],
})
export class StatisticsPage implements OnInit {

  @ViewChild('chartCanvas') chartCanvas : ElementRef;
  data : any = [];
  canvasChart : Chart;

  matchingstatistics: MatchingStat[] = [];
  connexionStat: string; emailSentStat: string; emailFailedStat: string;
  numUsers: boolean = false;
  numgroups = 0;
  numconves = 0;
  numfailed = 0;
  numMatchingSection = 0;

  constructor(private matchingstatService: MatchingStatService, public translate: TranslateService) { }

  ngOnInit() {
    this.connexionStat = this.translate.instant('NUM_CONNEXIONS');
    this.emailSentStat = this.translate.instant('EMAIL_SENT');
    this.emailFailedStat = this.translate.instant('EMAIL_FAILED');
    this.getMathingStats();
  }

  getMathingStats() {
    let numgroups = [];
    let numemails = [];
    let numfailed = [];
    let labels = [];

    this.matchingstatService.getMatchingStats().subscribe((matchingStats) => {
      if(matchingStats != null) {
        this.numMatchingSection = matchingStats.length;

        this.matchingstatistics = matchingStats
        for(const mathingStat of matchingStats) {
          labels.push(formatDate(mathingStat.createdAt, "MMM d, y", "en" /**navigator.language.split("-")[0] ||  */))
          numgroups.push(mathingStat.numGroups)
          numemails.push(mathingStat.numConversations)
          numfailed.push(mathingStat.numFailed)
          this.numconves += mathingStat.numConversations
          this.numfailed += mathingStat.numFailed 
        }

        this.data = {
          labels: labels,
          datasets: [{
            label: this.connexionStat,
            data: numgroups,
            backgroundColor: 'rgba(153, 102, 255)',
            borderColor: 'rgb(153, 102, 255)',
            borderWidth: 1
          },
          {
            label: this.emailSentStat,
            data: numemails,
            backgroundColor: 'rgba(75, 192, 192)',
            borderColor: 'rgb(75, 192, 192)',
            borderWidth: 1
          },
          {
            label: this.emailFailedStat,
            data: numfailed,
            backgroundColor: 'rgba(255, 99, 132)',
            borderColor: 'rgb(255, 99, 132)',
            borderWidth: 1
          }]
        };
        this.changeChart({detail: {
          value : 'bar'
        }});
      }
    });
  }

  changeChart( event: any ) {
    const type = event.detail.value || 'bar';
    if ( this.canvasChart ) {
      this.canvasChart.destroy();
    }
    this.canvasChart = new Chart(this.chartCanvas.nativeElement, {
      type: type,
      data: this.data,
      options: {
        indexAxis: 'x'
      }
    });
  }

  addRandom( points: any ) : number {
    return Number(points) - Number( Math.floor((Math.random() * 100) + 1) );
  }
}
