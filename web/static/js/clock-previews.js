// Pixel layouts transcribed from the user's numbered device photographs.
// Only clock-colored cells use the selected color; decoration stays fixed.
(function () {
  var small = {
    '0':['111','101','101','101','111'],
    '1':['010','110','010','010','111'],
    '2':['111','001','111','100','111'],
    '4':['101','101','111','001','001'],
    '5':['111','100','111','001','111'],
    '6':['111','100','111','101','111'],
    '7':['111','001','001','001','001'],
    '8':['111','101','111','101','111']
  };
  var large = {
    '1':['00100','01100','10100','00100','00100','00100','11111'],
    '9':['01110','10001','10001','01111','00001','10001','01110'],
    '5':['11111','10000','10000','11110','00001','10001','01110']
  };
  var rainbow = ['#ff5522','#ff263f','#ed30ad','#ba20e7','#7827e4','#4020ec','#242cff','#1764ff','#159dff','#19c6f2','#20dcc8','#62e52a','#a5ef23','#d9ed29','#ffd441','#ffa631'];
  function layout(style) {
    var cells = new Map();
    function dot(x,y,color) { cells.set(x+','+y,{x:x,y:y,color:color || 'clock'}); }
    function glyph(font,char,x,y,color) {
      font[char].forEach(function(row,dy) {
        Array.from(row).forEach(function(bit,dx) { if(bit==='1') dot(x+dx,y+dy,color); });
      });
    }
    function stacked(minutes) {
      glyph(small,'2',4,3); glyph(small,'0',9,3);
      glyph(small,minutes[0],4,10); glyph(small,minutes[1],9,10);
    }
    function horizontal(minutes,y,colon) {
      glyph(small,'2',0,y); glyph(small,'0',4,y);
      glyph(small,minutes[0],9,y); glyph(small,minutes[1],13,y);
      if(colon) { dot(7,y+1); dot(7,y+3); }
    }
    if(style===0) {
      glyph(large,'1',3,0); glyph(large,'9',9,0);
      glyph(large,'5',3,9); glyph(large,'9',9,9); dot(1,9); dot(1,13);
    } else if(style===1) {
      // Inverted, larger digits: the illuminated background is configurable.
      for(var y=0;y<16;y++) for(var x=0;x<16;x++) dot(x,y);
      var inverse = {
        '2':['01110','10001','00001','00010','00100','01000','11111'],
        '0':['01110','10001','10001','10001','10001','10001','01110'],
        '1':['00100','01100','00100','00100','00100','00100','11111']
      };
      [['2',4,1],['0',10,1],['0',4,9],['1',10,9]].forEach(function(v) {
        inverse[v[0]].forEach(function(row,dy) { Array.from(row).forEach(function(bit,dx) { if(bit==='1') cells.delete((v[1]+dx)+','+(v[2]+dy)); }); });
      });
    } else if(style>=2 && style<=6) {
      stacked(style===2?'02':style<=4?'04':'05');
      if(style===3) [0,1,14,15].forEach(function(x){dot(x,14,'#d6ceff');});
      if(style===4) [[0,0],[15,0],[0,15],[15,15]].forEach(function(p){dot(p[0],p[1],'#ed21eb');});
      if(style===5) {
        for(var i=0;i<16;i++) { dot(i,0,rainbow[(i+5)%16]); dot(15,i,rainbow[(i+1)%16]); dot(15-i,15,rainbow[(i+11)%16]); dot(0,15-i,rainbow[(i+9)%16]); }
      }
      if(style===6) [[0,2],[2,3],[9,2],[1,8],[7,8],[12,8],[14,8],[2,10],[14,10],[1,12],[15,12],[5,15],[14,15]].forEach(function(p){dot(p[0],p[1]);});
    } else if(style>=7 && style<=11) {
      horizontal(style<=8?'06':style===9?'07':'08',style<=8?3:6,style>=9);
      if(style===7) for(var j=1;j<16;j+=2) dot(j,14,'#e1d5ff');
      if(style===9) [[0,0],[0,1],[15,14],[15,15]].forEach(function(p){dot(p[0],p[1],'#ffd33e');});
      if(style===9) [[15,0],[15,1],[0,14],[0,15]].forEach(function(p){dot(p[0],p[1],'#8933ff');});
      if(style===10) [1,2,14,15].forEach(function(y){dot(7,y,'#ff243d');dot(8,y,'#ff243d');});
      if(style===11) {
        [[0,1],[1,1],[2,1],[3,1],[4,1],[5,1],[0,2],[1,2],[2,2],[3,2],[5,2],[0,3],[1,3],[2,3],[3,3],[5,3],[0,4],[2,4],[7,1],[7,2],[2,12],[0,13],[2,13],[0,14],[1,14],[2,14],[3,14],[5,14],[0,15],[1,15],[2,15],[3,15],[4,15],[5,15],[7,14],[7,15]].forEach(function(p){dot(p[0],p[1],p[1]===4 || p[1]===12?'#ff8431':'#ff3530');});
      }
    } else if(style===12 || style===13) {
      if(style===13) for(var cy=0;cy<16;cy++) for(var cx=0;cx<16;cx++) dot(cx,cy,'#168caa');
      else {
        for(var k=2;k<14;k++) {dot(k,0,'#aeb4ee');dot(k,15,'#aeb4ee');dot(0,k,'#aeb4ee');dot(15,k,'#aeb4ee');}
        [[1,1],[14,1],[1,14],[14,14]].forEach(function(p){dot(p[0],p[1],'#aeb4ee');});
      }
      [[7,0],[8,0],[0,7],[0,8],[15,7],[15,8],[7,15],[8,15]].forEach(function(p){dot(p[0],p[1]);});
      // Style 13 shows 20:30: the short hand is between eight and nine.
      var hand=style===12?[[8,3],[7,4],[8,4],[6,5],[7,5],[6,6],[7,6],[6,7],[6,8]]:[[7,8],[5,9],[6,9],[7,9],[5,10]];
      hand.forEach(function(p){dot(p[0],p[1],'#e7dcff');});
      // Style 13's cyan minute hand points straight down at :30.
      // Its fixed color is brighter than the cyan face, not clock-colored.
      if(style===13) for(var minuteY=8;minuteY<=13;minuteY++) dot(8,minuteY,'#26dfd6');
      if(style===12) [[5,9],[6,9],[4,10],[5,10],[4,11]].forEach(function(p){dot(p[0],p[1],'#ffbc38');});
    } else if(style===15) {
      horizontal('10',6,true);
      for(var r=0;r<16;r++){dot(r,4,rainbow[r]);dot(r,12,rainbow[15-r]);}
    }
    return Array.from(cells.values());
  }
  function render() {
    var picker = document.getElementById('clock-color');
    if(!picker) return;
    document.querySelectorAll('.clock-preview').forEach(function(svg) {
      svg.style.setProperty('--clock-color',picker.value);
      if(svg.childElementCount) return;
      layout(Number(svg.dataset.style)).forEach(function(cell) {
        var rect=document.createElementNS('http://www.w3.org/2000/svg','rect');
        rect.setAttribute('x',cell.x+0.08);rect.setAttribute('y',cell.y+0.08);
        rect.setAttribute('width','0.84');rect.setAttribute('height','0.84');rect.setAttribute('rx','0.06');
        rect.setAttribute('fill',cell.color==='clock'?'var(--clock-color)':cell.color);
        svg.appendChild(rect);
      });
    });
  }
  document.addEventListener('input',function(event){if(event.target.id==='clock-color') render();});
  document.addEventListener('htmx:afterSwap',function(){render();});
  render();
})();
