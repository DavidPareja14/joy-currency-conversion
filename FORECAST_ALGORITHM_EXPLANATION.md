# 📊 Forecast Algorithm - Detailed Explanation

## 🎯 **Why 5 Days Instead of 30?**

- **API Limitations**: Each day requires a separate API call
- **30 days = 30 API calls** → Too expensive and slow
- **5 days = 5 API calls** → Much more efficient
- **Still provides meaningful trend analysis**

---

## 📈 **Step-by-Step Algorithm Explanation**

### **Step 1: Data Collection (5 Days)**
```
Example: EUR/USD rates over 5 days
Day 1: 1.0800
Day 2: 1.0820  
Day 3: 1.0810
Day 4: 1.0830
Day 5: 1.0840
```

### **Step 2: Calculate Average**
```
Average = (1.0800 + 1.0820 + 1.0810 + 1.0830 + 1.0840) / 5
Average = 5.4100 / 5 = 1.0820
```

### **Step 3: Calculate Standard Deviation**

**What is Standard Deviation?**
- Measures how "scattered" or "spread out" the data is
- **Low standard deviation** = rates are stable and consistent
- **High standard deviation** = rates are volatile and unpredictable

**Calculation:**
```
For each rate, calculate: (rate - average)²

Day 1: (1.0800 - 1.0820)² = (-0.0020)² = 0.000004
Day 2: (1.0820 - 1.0820)² = (0.0000)² = 0.000000
Day 3: (1.0810 - 1.0820)² = (-0.0010)² = 0.000001
Day 4: (1.0830 - 1.0820)² = (0.0010)² = 0.000001
Day 5: (1.0840 - 1.0820)² = (0.0020)² = 0.000004

Variance = (0.000004 + 0.000000 + 0.000001 + 0.000001 + 0.000004) / 5
Variance = 0.000010 / 5 = 0.000002

Standard Deviation = √0.000002 = 0.0014
```

**Why is this important?**
- **Low standard deviation** → High confidence in prediction
- **High standard deviation** → Low confidence in prediction

### **Step 4: Trend Analysis**

**Compare first half vs second half:**
```
First half (Days 1-2): Average = (1.0800 + 1.0820) / 2 = 1.0810
Second half (Days 4-5): Average = (1.0830 + 1.0840) / 2 = 1.0835

Trend = (1.0835 - 1.0810) / 1.0810 = 0.0023 (0.23% increase)
```

**What does this mean?**
- **Positive trend** → Rate is increasing
- **Negative trend** → Rate is decreasing
- **Zero trend** → Rate is stable

### **Step 5: Predict Next Day**

**Conservative Prediction Formula:**
```
Predicted Rate = Average + (Trend × 0.5)

Predicted Rate = 1.0820 + (0.0023 × 0.5)
Predicted Rate = 1.0820 + 0.00115
Predicted Rate = 1.08315
```

**Why 50% of trend?**
- **100% trend** → Too aggressive, might be wrong
- **0% trend** → Too conservative, ignores trends
- **50% trend** → Balanced approach

### **Step 6: Calculate Confidence**

**Base confidence: 50%**

**Coefficient of Variation:**
```
Coefficient of Variation = Standard Deviation / Average
Coefficient of Variation = 0.0014 / 1.0820 = 0.0013 (0.13%)
```

**Confidence adjustments:**
- **Coefficient < 5%** → +30% confidence (very stable)
- **Coefficient < 10%** → +20% confidence (stable)
- **Coefficient < 20%** → +10% confidence (somewhat stable)
- **5 data points** → +10% confidence
- **4 data points** → +5% confidence

**Final confidence calculation:**
```
Base: 50%
Stability bonus: +30% (because 0.13% < 5%)
Data points bonus: +10% (because we have 5 data points)
Final confidence: 50% + 30% + 10% = 90%
```

**Confidence range: 30% - 90%**

---

## 🔍 **Real Example with Different Scenarios**

### **Scenario 1: Stable Currency (High Confidence)**
```
Rates: [1.0800, 1.0805, 1.0802, 1.0803, 1.0801]
Average: 1.0802
Standard Deviation: 0.0002 (very low)
Trend: 0.0001 (almost no trend)
Prediction: 1.0802 (very close to average)
Confidence: 90% (high confidence because very stable)
```

### **Scenario 2: Volatile Currency (Low Confidence)**
```
Rates: [1.0800, 1.0900, 1.0700, 1.0850, 1.0750]
Average: 1.0800
Standard Deviation: 0.0071 (high)
Trend: -0.0015 (decreasing)
Prediction: 1.0793 (slightly below average)
Confidence: 30% (low confidence because very volatile)
```

### **Scenario 3: Trending Currency (Medium Confidence)**
```
Rates: [1.0800, 1.0820, 1.0840, 1.0860, 1.0880]
Average: 1.0840
Standard Deviation: 0.0032 (medium)
Trend: 0.0074 (strong upward trend)
Prediction: 1.0877 (above average due to trend)
Confidence: 60% (medium confidence)
```

---

## 🎯 **Why This Algorithm Works**

1. **Uses Real Data**: Based on actual historical rates, not random numbers
2. **Conservative Approach**: 50% trend adjustment prevents extreme predictions
3. **Confidence Indicator**: Tells you how reliable the prediction is
4. **API Efficient**: Only 5 API calls instead of 30
5. **Statistically Sound**: Uses proper statistical methods

---

## 🚀 **API Usage**

```bash
# Test the forecast endpoint
curl "http://localhost:8080/api/v1/forecast?origin=EUR&destination=USD"
```

**Response:**
```json
{
  "origin": {"code": "EUR", "country": "Eurozone"},
  "destination": {"code": "USD", "country": "United States"},
  "predicted_date": "2025-01-15",
  "predicted_rate": 1.08315,
  "confidence": 0.90,
  "last_30_days": {
    "average": 1.0820
  },
  "timestamp": "2025-01-14T12:00:00Z",
  "rates_source": "api.exchangeratesapi.io"
}
```

---

## ⚠️ **Important Notes**

- **Minimum 3 days** of data required
- **Confidence range**: 30% - 90%
- **Conservative prediction**: Uses 50% of trend + 50% of average
- **API efficient**: Only 5 API calls maximum
- **Real-time data**: Uses actual exchange rates from APIs
