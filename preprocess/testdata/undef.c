#define A 10 + __LINE__ + B
#define B 20 + A + C
#define D B
#undef B
#define C 30 + B
#define E int a = 10;
#undef E

int main() {
    E
    printf("%d %d %d\n", D, C);
}