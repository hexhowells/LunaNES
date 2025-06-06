<p align="center">
  <img src="https://github.com/hexhowells/LunaNES/blob/main/logo.jpg" width=40%>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white"/>
  <img src="https://img.shields.io/github/license/hexhowells/LunaNES"/>
  <img src="https://custom-icon-badges.demolab.com/badge/Windows-0078D6?logo=windows11&logoColor=white"/>
  <img src="https://img.shields.io/badge/Linux-FCC624?logo=linux&logoColor=black"/>
  <img src="https://img.shields.io/badge/macOS-000000?logo=apple&logoColor=F0F0F0"/>
</p>

LunaNES is an NES emulator written in go. It is fully functional with mapper 0 with more mappers soon to be developed. Currently the emulator does not support sound. The current configuration of the emulator recieves input from a USB NES controller, the controllers VID and PID will be needed to ensure LunaNES connects to the correct device.

---

### Game Screenshots
<p align="center">
  <img src="https://github.com/hexhowells/LunaNES/blob/main/games.png" width=85%>
</p>

### Currently supported mappers
- [Mapper 0](https://nesdir.github.io/mapper0.html)
- [Mapper 1](https://nesdir.github.io/mapper1.html) (planned)
- [Mapper 2](https://nesdir.github.io/mapper2.html) (works with some games)
- [Mapper 3](https://nesdir.github.io/mapper3.html) (planned)
- [Mapper 4](https://nesdir.github.io/mapper4.html) (kinda works)

### Todo
- [ ] Support more mappers
- [ ] Optimise PPU clock function (currently a little slow with too many sprites on screen)
- [ ] Support 2 controllers
- [ ] Support keyboard input
- [ ] Implement sound

