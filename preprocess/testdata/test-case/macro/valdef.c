#define A 10 + __LINE__
#define B 20 + A
#define C 30 + B

int main() {
    printf("%d\n", C);
}